package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"os"
)

const (
	rootPath    = "/home/www/dream/prodapproot/modules"
	xmlFileName = "config.xml"
	dirSep      = "/"
	DSNpsold      = "ps:pass_to_ps@/ps?charset=utf8"
	DSNold      = "dlprod:pass_to_dlprod@/dlprod?charset=utf8"
)

type TableModules struct {
	IdModule int
	Name     string
	Active   int
	Version  string
}

type Modules struct {
	Id                int
	IdNew             int
	IdOld             int
	PathnameOld       string
	PathnameNew       string
	NameOld           string
	NameNew           string
	AuthorOld         string
	AuthorNew         string
	VersionOld        string
	VersionNew        string
	ActiveOld         int
	ActiveNew         int
	IsConfigurableOld int
	IsConfigurableNew int
	AvailableUrl      string
	DescriptionOld    string
	DescriptionNew    string
}

func main() {
	dbps, err := sql.Open("mysql", DSNpsold)
	checkErr(err)
	defer dbps.Close()

	//Connection to psnew database:
	dbold, err := sql.Open("mysql", DSNold)
	checkErr(err)
	defer dbold.Close()
	dir, err := ioutil.ReadDir(rootPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, fn := range dir {
		//If is a directory then done normaly
		if !fn.Mode().IsDir() {
			continue
		}
		discoverPath := rootPath + dirSep + fn.Name() + dirSep + xmlFileName
		bs, err := ioutil.ReadFile(discoverPath)

		//If file exists then done normaly
		if os.IsNotExist(err) {
			continue
		}
		str := string(bs)
		module := GetModulesDataFromXml(str)
		fmt.Println(discoverPath, ":", module.DisplayName, "(", module.Author, ")")
		//var item TableModules
		item := TableModules{}
		record := dbold.QueryRow("SELECT id_module, name, active, version FROM dlprod.ps_module WHERE name=?", module.Pathname)
		err = record.Scan(&item.IdModule, &item.Name, &item.Active, &item.Version)
		switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			//return
		case nil:
			fmt.Println("Record:", item)
		default:
			panic(err)
		}
		updateGatheredFromXmlAndDb(dbps, module, item)
	}
}

func updateGatheredFromXmlAndDb(dbd *sql.DB, m XmlModule,r TableModules) {
	// update from xml file:
    //Check record is here:
    item := TableModules{}
    record := dbd.QueryRow("SELECT name_old, active_old, version_old FROM ps.modules WHERE pathname_old=?", m.Pathname)
    err := record.Scan(&item.Name, &item.Active, &item.Version)
    switch err {
    case sql.ErrNoRows:
        fmt.Println("No rows were returned in ps.modules")
    case nil:
        fmt.Println("Record in ps.modules founded:", item)
    default:
        panic(err)
    }

    // If no such record:
    if err == sql.ErrNoRows {
        stmt, err := dbd.Prepare("INSERT ps.modules SET id_old=?, pathname_old=?, name_old=?, author_old=?, version_old=?, active_old=?, description_old=?, is_configurable_old=?")
        checkErr(err)
        res, err := stmt.Exec(
            r.IdModule,
            m.Pathname,
            m.DisplayName,
            m.Author,
            m.Version,
            r.Active,
            m.Description,
            m.IsConfigurable,
        )
        checkErr(err)
        affect, err := res.RowsAffected()
        checkErr(err)
        fmt.Println(m.Pathname, m.DisplayName, affect, "- is inserted in DB")
    } else {
        stmt, err := dbd.Prepare("UPDATE ps.modules SET id_old=?, pathname_old=?, name_old=?, author_old=?, version_old=?, active_old=?, description_old=?, is_configurable_old=? WHERE pathname_old = ?")
        checkErr(err)
        res, err := stmt.Exec(
            r.IdModule,
            m.Pathname,
            m.DisplayName,
            m.Author,
            m.Version,
            r.Active,
            m.Description,
            m.IsConfigurable,
            m.Pathname,
        )
        checkErr(err)
        affect, err := res.RowsAffected()
        checkErr(err)
        fmt.Println(m.Pathname, m.DisplayName, affect, "- is updated in DB")
    }
    return
}

type XmlModule struct {
	XMLName          xml.Name `xml:"module"`
	Pathname         string   `xml:"name"`
	DisplayName      string   `xml:"displayName"`
	Version          string   `xml:"version"`
	Description      string   `xml:"description"`
	Author           string   `xml:"author"`
	Tab              string   `xml:"tab"`
	IsConfigurable   string   `xml:"is_configurable"`
	NeedInstance     string   `xml:"need_instance"`
	LimitedCountries string   `xml:"limited_countries"`
}

func GetModulesDataFromXml(rawXMLdata string) XmlModule {
	var data XmlModule
	xml.Unmarshal([]byte(rawXMLdata), &data)
	return data
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
