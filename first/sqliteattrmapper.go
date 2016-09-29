// SQLite attribute mapper

package first

import (
	"log"
	"database/sql"
	"fmt"
	"bytes"
	"strings"
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

func (attrMapper *SQLiteAttrMapper) ReadEntries(query string) []FileMetadataEntry {
	rows, err := attrMapper.db.Query(query)
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

func (attrMapper *SQLiteAttrMapper) ReadEntries2(query string) []FileMetadataEntry {
	rows, err := attrMapper.db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []FileMetadataEntry
	for rows.Next() {
		entry := FileMetadataEntry{}
		err2 := rows.Scan(&entry.fileID)
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

func (attrMapper *SQLiteAttrMapper) ReadEntry() []FileMetadataEntry {
	sql_readAll := `
	SELECT fileID, attribute, value FROM FileMetadata
	`
	return attrMapper.ReadEntries(sql_readAll)
}

func (attrMapper *SQLiteAttrMapper) QueryBuilderForMultipleUUIDSelections(attributes *QueryKeyValue) (string, bool) {
	_, mainQuery, foundQuery := QueryBuilderForUUIDSelection(attributes)
	if !foundQuery {
		return "", false
	}
	preResults := attrMapper.ReadEntries2(mainQuery)
	log.Println(preResults)
	return fmt.Sprintf("SELECT fileID, attribute, value FROM FileMetadata WHERE fileID IN ( %v )", mainQuery), true
}

func (attrMapper *SQLiteAttrMapper) FindAllMatchingQueries(attributes *QueryKeyValue) ([]UUIDToQuery, bool) {
	builtQuery, querySuccess := attrMapper.QueryBuilderForMultipleUUIDSelections(attributes)
	if querySuccess {
		log.Println("Built this query for all matching uuids ", builtQuery)
		results := attrMapper.ReadEntries(builtQuery)
		log.Println(results)
		if len(results) > 0 {

			uuidToAttributeValue := make(map[string]map[string]string, 0)
			for _, entry := range results {
				AddKeyValuePairToUUIDMap(entry.attribute, entry.value, entry.fileID, uuidToAttributeValue)
			}

			queryKeyValues := []UUIDToQuery{}
			for uuid := range uuidToAttributeValue {
				queryKeyValue := QueryKeyValue{
					uuidToAttributeValue[uuid],
				}
				foundQuery := UUIDToQuery{
					uuid,
					queryKeyValue,
				}
				queryKeyValues = append(queryKeyValues, foundQuery)
			}
			return queryKeyValues, true
		}
	}
	return nil, false
}


func QueryBuilderForUUIDSelection(attributes *QueryKeyValue) (string, string, bool) {
	if attributes == nil || attributes.keyValue == nil || len(attributes.keyValue) == 0 {
		return "", "", false
	}
	var queryBuf bytes.Buffer
	var attrBuf bytes.Buffer
	for key, value := range attributes.keyValue {
		if queryBuf.Len() != 0 && attrBuf.Len() != 0 {
			queryBuf.WriteString("INTERSECT ")
			attrBuf.WriteString(" AND ")
		}
		queryBuf.WriteString(fmt.Sprintf("SELECT fileID FROM FileMetadata WHERE attribute='%v' AND value='%v' ", key, value))
		attrBuf.WriteString(fmt.Sprintf("attribute!='%v'", key))
		someStr := queryBuf.String()
		log.Println(someStr)
	}
	secondary := queryBuf.String()
	if attrBuf.Len() != 0 {
		queryBuf.WriteString("EXCEPT SELECT fileID FROM FileMetadata WHERE ")
		queryBuf.Write(attrBuf.Bytes())
	}
	return queryBuf.String(), secondary, true
}

func (attrMapper *SQLiteAttrMapper) GetAddedUUID(attributes *QueryKeyValue, isFile bool) (string, bool) {
	log.Println("Reading all entries")
	log.Println(attrMapper.ReadEntry())

	builtQuery, secondary, querySuccess := QueryBuilderForUUIDSelection(attributes)
	if querySuccess {
		log.Println("Built this query ", builtQuery)
		results := attrMapper.ReadEntries2(builtQuery)
		log.Println(results)
		if len(results) == 0 && !isFile { //Definitely not a file, potentially could be a directory
			if strings.EqualFold(builtQuery, secondary) {
				return "", false
			}
			results = attrMapper.ReadEntries2(secondary)
		}
		if !isFile && len(results) > 0 {
			return "", true //It's a file or a directory
		}
		if len(results) < 2 {
			for _, result := range results {
				return result.fileID, true
			}
		} else {
			log.Fatal("Found ", len(results), " results instead of 1")
		}
	}
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

	uuidStr, attributeAdded := attrMapper.GetAddedUUID(attributes, true)
	if attributeAdded {
		return uuidStr
	}
	return CreateNewUUID(attributes, attrMapper.AddQueryToUUID)

}
