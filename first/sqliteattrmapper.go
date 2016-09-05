// In memory attribute mapper

package first

import "log"

type SQLiteAttrMapper struct {
	AttrMapper
}

func NewSQLiteAttrMapper() *SQLiteAttrMapper {
	sqliteAttrMapper := &SQLiteAttrMapper{
	}
	return sqliteAttrMapper
}

func (attrMapper *SQLiteAttrMapper) GetAddedUUID(attributes *QueryKeyValue) (string, bool) {
	return "", false
}

func (attrMapper *SQLiteAttrMapper) CreateFromQuery(attributes *QueryKeyValue) string {
	log.Println("Not implemented")
	return "foo"
}
