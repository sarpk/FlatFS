// MySQL attribute mapper

package FlatFS

import (
	"log"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type MySQLAttrMapper struct {
	AttrMapper
	db *sql.DB
}

//type FileMetadataEntry struct {
//	fileID    string
//	attribute string
//	value     string
//}

func NewMySQLAttrMapper() *MySQLAttrMapper {

	const dataSource = "root:changeit@/flatfs"

	mySqlAttrMapper := &MySQLAttrMapper{
		db: InitMySQLeDB(dataSource),
	}
	mySqlAttrMapper.CreateTable()
	return mySqlAttrMapper
}

func (attrMapper *MySQLAttrMapper) Close() {
	attrMapper.db.Close()
}

func (attrMapper *MySQLAttrMapper) CreateTable() {
	// create table if not exists
	sql_table := `
	CREATE TABLE IF NOT EXISTS FileMetadata(
		fileID TEXT NOT NULL,
		attribute TEXT NOT NULL,
		value TEXT,
		PRIMARY KEY (fileID(128), attribute(128))
	);
	`

	_, err := attrMapper.db.Exec(sql_table)
	if err != nil {
		panic(err)
	}
}

func (attrMapper *MySQLAttrMapper) ReadEntries(query string) []FileMetadataEntry {
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

func (attrMapper *MySQLAttrMapper) ReadEntries2(query string) []FileMetadataEntry {
	start := time.Now()
	rows, err := attrMapper.db.Query(query)
	log.Println("Query is ", query)
	log.Println("Query took ", time.Since(start))
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []FileMetadataEntry
	log.Println("FileMetadataEntry init took ", time.Since(start))

	for rows.Next() {
		log.Println("rows iteration took ", time.Since(start))
		entry := FileMetadataEntry{}
		err2 := rows.Scan(&entry.fileID)
		log.Println("rows.Scan  took ", time.Since(start))
		if err2 != nil {
			panic(err2)
		}
		result = append(result, entry)
		log.Println("append  took ", time.Since(start))
	}
	return result
}

func (attrMapper *MySQLAttrMapper) StoreEntry(entries []FileMetadataEntry) {
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

func (attrMapper *MySQLAttrMapper) ReadEntry() []FileMetadataEntry {
	sql_readAll := `
	SELECT fileID, attribute, value FROM FileMetadata
	`
	return attrMapper.ReadEntries(sql_readAll)
}

func (attrMapper *MySQLAttrMapper) QueryBuilderForMultipleUUIDSelections(attributes *QueryKeyValue) (string, bool) {
	_, mainQuery, foundQuery := QueryBuilderForUUIDSelection(attributes)
	if !foundQuery {
		return "", false
	}
	//preResults := attrMapper.ReadEntries2(mainQuery)
	//log.Println(preResults)
	return fmt.Sprintf("SELECT fileID, attribute, value FROM FileMetadata WHERE fileID IN ( %v )", mainQuery), true
}

func (attrMapper *MySQLAttrMapper) FindAllMatchingQueries(attributes *QueryKeyValue) ([]UUIDToQuery, bool) {
	builtQuery, querySuccess := attrMapper.QueryBuilderForMultipleUUIDSelections(attributes)
	if querySuccess {
		//log.Println("Built this query for all matching uuids ", builtQuery)
		results := attrMapper.ReadEntries(builtQuery)
		//log.Println(results)
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


//func QueryBuilderForUUIDSelection(attributes *QueryKeyValue) (string, string, bool) {
//	if attributes == nil || attributes.keyValue == nil || len(attributes.keyValue) == 0 {
//		return "", "", false
//	}
//	var queryBuf bytes.Buffer
//	//var attrBuf bytes.Buffer
//	for key, value := range attributes.keyValue {
//		if queryBuf.Len() != 0 {
//			queryBuf.WriteString("INTERSECT ")
//			//attrBuf.WriteString(" AND ")
//		}
//		queryBuf.WriteString(fmt.Sprintf("SELECT fileID FROM FileMetadata WHERE attribute='%v' AND value='%v' ", key, value))
//		//attrBuf.WriteString(fmt.Sprintf("attribute!='%v'", key))
//		//someStr := queryBuf.String()
//		//log.Println(someStr)
//	}
//	secondary := queryBuf.String()
//	//if attrBuf.Len() != 0 {
//	//	queryBuf.WriteString("EXCEPT SELECT fileID FROM FileMetadata WHERE ")
//	//	queryBuf.Write(attrBuf.Bytes())
//	//}
//	return queryBuf.String(), secondary, true
//}

func (attrMapper *MySQLAttrMapper) DeleteUUIDFromQuery(attributes *QueryKeyValue, uuid string) {
	for key, value := range attributes.keyValue {
		deleteQuery := fmt.Sprintf("Delete FROM FileMetadata WHERE fileID='%v' AND attribute='%v' AND value='%v'", uuid, key, value)
		attrMapper.ReadEntries2(deleteQuery)
	}
}

func (attrMapper *MySQLAttrMapper) GetAddedUUID(attributes *QueryKeyValue, queryType QueryType) (string, bool) {
	//log.Println("Reading all entries")
	//log.Println(attrMapper.ReadEntry())
	//start := time.Now()

	builtQuery, secondary, querySuccess := QueryBuilderForUUIDSelection(attributes)
	//log.Println("QueryBuilderForUUIDSelection took ", time.Since(start))
	if querySuccess {
		//log.Println("Built this query ", builtQuery)
		results := attrMapper.ReadEntries2(builtQuery)
		//log.Println("ReadEntries2 took ", time.Since(start))
		//log.Println(results)
		if len(results) == 0 && !queryType.fileSpec { //Definitely not a file, potentially could be a directory
			if strings.EqualFold(builtQuery, secondary) {
				return "", false
			}
			results = attrMapper.ReadEntries2(secondary)
			//log.Println("ReadEntries2 secondary took ", time.Since(start))
		}
		if !queryType.fileSpec && len(results) > 0 {
			//log.Println("!queryType.fileSpec took ", time.Since(start))
			return "", true //It's a file or a directory
		}
		if len(results) < 2 {
			for _, result := range results {
				//log.Println("results iteration took ", time.Since(start))
				return result.fileID, true
			}
		} else {
			log.Fatal("Found ", len(results), " results instead of 1")
		}
	}
	return "", false
}

func (attrMapper *MySQLAttrMapper) AddQueryToUUID(key, value, uuid string) {
	file := FileMetadataEntry{
		fileID: uuid,
		attribute: key,
		value: value,
	}
	attrMapper.StoreEntry([]FileMetadataEntry{file})
}

func (attrMapper *MySQLAttrMapper) CreateFromQuery(attributes *QueryKeyValue) string {
	//log.Println("Not implemented")
	//start := time.Now()
	fileSpecQueryType := createFileSpecQueryType()
	//log.Println("fileSpecQueryType took ", time.Since(start))
	uuidStr, attributeAdded := attrMapper.GetAddedUUID(attributes, fileSpecQueryType)
	//log.Println("GetAddedUUID took ", time.Since(start))
	if attributeAdded {
		//log.Println("GetAddedUUID took ", time.Since(start))
		return uuidStr
	}
	//log.Println("GetAddedUUID took ", time.Since(start))
	result := CreateNewUUID(attributes, attrMapper.AddQueryToUUID)
	//log.Println("CreateNewUUID took ", time.Since(start))
	return result

}
