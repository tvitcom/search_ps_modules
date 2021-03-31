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
	rootPath    = "/home/www/ps/newapproot/modules"
	xmlFileName = "config.xml"
	dirSep      = "/"
	DSNpsold      = "ps:pass_to_ps@/ps?charset=utf8"
	DSNpsnew    = "psnew:pass_to_psnew@/psnew?charset=utf8"
)

type TableModules struct {
	IdModule       int
	Name           string
	Active         int
	Version        string
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
    //Connection to old database:
	dbold, err := sql.Open("mysql", DSNpsold)
    checkErr(err)
	defer dbold.Close()

    //Connection to psnew database:
	dbnew, err := sql.Open("mysql", DSNpsnew)
    checkErr(err)
	defer dbnew.Close()
    //Test connections
    err = dbold.Ping()
    checkErr(err)
    err = dbnew.Ping()
    checkErr(err)

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
		outp := GetModulesDataFromXml(str)
		fmt.Println(discoverPath, ":", outp.DisplayName, "(", outp.Author, ")")
        //var item TableModules
        item := TableModules{}
        record := dbnew.QueryRow("SELECT id_module, name, active, version FROM psnew.ps_module WHERE name=?", outp.Pathname)
        err = record.Scan(&item.IdModule,&item.Name, &item.Active, &item.Version)
        switch err {
            case sql.ErrNoRows:
                fmt.Println("No rows were returned!")
                //return
            case nil:
                fmt.Println("Record:", item)
            default:
              panic(err)
		}
        updateGatheredFromXmlAndDb(dbold, outp, item)
	}
}

func updateGatheredFromXmlAndDb(dbd *sql.DB, m XmlModule, r TableModules) {
	// update from xml file:
    //Check record is here:
    item := TableModules{}
    record := dbd.QueryRow("SELECT id, name_new, active_new, version_new FROM ps.modules WHERE pathname_new=? OR pathname_old = ? OR description_old = ?", m.Pathname, m.Pathname, m.Description)
    err := record.Scan(&item.IdModule, &item.Name, &item.Active, &item.Version)
    switch err {
    case sql.ErrNoRows:
        fmt.Println("No rows were returned in ps.modules")
        //return
    case nil:
        fmt.Println("Record in ps.modules founded:Id=", item.IdModule)
    default:
        panic(err)
    }

    // If no such record:
    if err == sql.ErrNoRows {
        stmt, err := dbd.Prepare("INSERT ps.modules SET id_new=?, pathname_new=?, name_new=?, author_new=?, version_new=?, active_new=?, is_configurable_new=?, description_new=?")
        checkErr(err)
        res, err := stmt.Exec(
            r.IdModule,
            m.Pathname,
            m.DisplayName,
            m.Author,
            m.Version,
            r.Active,
            m.IsConfigurable,
            m.Description,
        )
        checkErr(err)
        affect, err := res.RowsAffected()
        checkErr(err)
        fmt.Println(m.Pathname, m.DisplayName, affect, "- is inserted in DB")
    } else {
    // If has record:
        stmt, err := dbd.Prepare("UPDATE ps.modules SET id_new=?, pathname_new=?, name_new=?, author_new=?, version_new=?, active_new=?, is_configurable_new=?, description_new=? WHERE id = ?")
        checkErr(err)
        res, err := stmt.Exec(
            r.IdModule,
            m.Pathname,
            m.DisplayName,
            m.Author,
            m.Version,
            r.Active,
            m.IsConfigurable,
            m.Description,
            item.IdModule,
        )
        checkErr(err)
        affect, err := res.RowsAffected()
        checkErr(err)
        fmt.Println(m.Pathname, m.DisplayName, affect, "- is updated in DB")
    }
    return
}

/* Example the Prestashop module xml description:
    rawData := `
<?xml version="1.0" encoding="UTF-8" ?>
<module>
    <name>belvg_related_products</name>
    <displayName><![CDATA[Related Products]]></displayName>
    <version><![CDATA[1.0.0]]></version>
    <description><![CDATA[Related Products]]></description>
    <author><![CDATA[BelVG]]></author>
    <tab><![CDATA[advertising_marketing]]></tab>
    <is_configurable>1</is_configurable>
    <need_instance>0</need_instance>
    <limited_countries></limited_countries>
</module>
`
*/

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
	//func GetModulesName(rawXMLdata string) {

	var data XmlModule
	xml.Unmarshal([]byte(rawXMLdata), &data)
	//fmt.Println(string(data.Name))
	return data
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
