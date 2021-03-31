It utils for gathering and prepare information about using Prestashop modules for currennt and new release. 

1. Firstly go run old_modules.go for gathering current xml modules descriptions of current Prestashop codebase.

2. Thirdly go run new_modules.go for gathering modules from xml descritions.

3. Then finally run next queries for representation modules states and need report creation.

```sql

-- modules-all:
CREATE table all_modules AS
SELECT /*id, id_old,id_new,*/pathname_old, name_new, author_old, pathname_new, name_old, author_new, available_url available, concat(left(description_old,28),"...") description_old, left(description_new,28) description_new
FROM modules
WHERE active_old > 0;

```

```sql

-- modules_done:57
CREATE table done_modules AS
SELECT /*id, id_old,id_new,*/ name_old, name_new,pathname_old, pathname_new, author_old, author_new, available_url available, concat(left(description_old,28),"...") description_old, left(description_new,28) description_new
FROM modules
WHERE active_old > 0 AND (BIT_LENGTH(pathname_new) > 0 OR BIT_LENGTH(available_url) > 0);

```

```sql

-- modules-discontinued:
CREATE table discontinued_modules AS
SELECT /*id, id_old,id_new,*/name_old, name_new, pathname_old, pathname_new, author_old, author_new, available_url available, is_configurable_old,is_configurable_new,available_url, concat(left(description_old,28),"...") description_old, left(description_new,28) description_new
FROM modules
WHERE active_old > 0 AND (BIT_LENGTH(pathname_new) = 0 AND BIT_LENGTH(available_url) = 0);

```
