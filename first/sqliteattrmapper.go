// In memory attribute mapper

package first

import (
	"log"
	"database/sql"
)

type SQLiteAttrMapper struct {
	AttrMapper
	db *sql.DB
}

type FileMetadataEntry struct {
	fileID    string
	attribute string
	value     string
}

func NewSQLiteAttrMapper() *SQLiteAttrMapper {

	const dbpath = "file_metadata.db"

	sqliteAttrMapper := &SQLiteAttrMapper{
		db: InitDB(dbpath),
	}
	sqliteAttrMapper.CreateTable()
	return sqliteAttrMapper
}

func (attrMapper *SQLiteAttrMapper) Close() {
	attrMapper.db.Close()
}

func (attrMapper *SQLiteAttrMapper) CreateTable() {
	// create table if not exists
	sql_table := `
	CREATE TABLE IF NOT EXISTS FileMetadata(
		fileID TEXT NOT NULL,
		attribute TEXT NOT NULL,
		value TEXT,
		PRIMARY KEY (fileID, attribute)
	);
	`

	_, err := attrMapper.db.Exec(sql_table)
	if err != nil {
		panic(err)
	}
}

func (attrMapper *SQLiteAttrMapper) ReadEntry() []FileMetadataEntry {
	sql_readAll := `
	SELECT fileID, attribute, value FROM FileMetadata
	`

	rows, err := attrMapper.db.Query(sql_readAll)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []FileMetadataEntry
	for rows.Next() {
		entry := FileMetadataEntry{}
		err2 := rows.Scan(&entry.fileID, &entry.attribute, &entry.value)
		if err2 != nil {
			panic(err2)
		}
		result = append(result, entry)
	}
	return result
}

func (attrMapper *SQLiteAttrMapper) StoreEntry(entries []FileMetadataEntry) {
	sql_addEntry := `
	INSERT OR REPLACE INTO FileMetadata(
		fileID,
		attribute,
		value
	) values(?, ?, ?)
	`

	stmt, err := attrMapper.db.Prepare(sql_addEntry)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	for _, entry := range entries {
		_, err2 := stmt.Exec(entry.fileID, entry.attribute, entry.value)
		if err2 != nil {
			panic(err2)
		}
	}
}

func (attrMapper *SQLiteAttrMapper) GetAddedUUID(attributes *QueryKeyValue) (string, bool) {
	log.Println("Reading all entries")
	log.Println(attrMapper.ReadEntry())
	return "", false
}

func (attrMapper *SQLiteAttrMapper) AddQueryToUUID(key, value, uuid string) {
	file := FileMetadataEntry{
		fileID: uuid,
		attribute: key,
		value: value,
	}
	attrMapper.StoreEntry([]FileMetadataEntry{file})
}

func (attrMapper *SQLiteAttrMapper) CreateFromQuery(attributes *QueryKeyValue) string {
	log.Println("Not implemented")

	uuidStr, attributeAdded := attrMapper.GetAddedUUID(attributes)
	if attributeAdded {
		return uuidStr
	}
	return CreateNewUUID(attributes, attrMapper.AddQueryToUUID)

}
