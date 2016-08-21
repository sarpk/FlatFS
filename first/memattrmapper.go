// In memory attribute mapper

package first

import (
	"log"
	"github.com/nu7hatch/gouuid"
)

type MemAttrMapper struct {
	AttrMapper
	queryToUuid map[string]map[string]string
}

func NewMemAttrMapper() *MemAttrMapper {
	memAttrMapper := &MemAttrMapper{
		queryToUuid: make(map[string]map[string]string, 0),
	}
	return memAttrMapper
}

func (attrMapper *MemAttrMapper) AddQueryToUUID(key, value, uuid string) {
	if attrMapper.queryToUuid[key] == nil {
		attrMapper.queryToUuid[key] = make(map[string]string, 0)
	}
	attrMapper.queryToUuid[key][value] = uuid
}


func (attrMapper *MemAttrMapper) CreateFromQuery(attributes *QueryKeyValue) string {
	log.Println("Mocking middleware")
	if uuid, err := uuid.NewV4(); err == nil {
		uuidStr := uuid.String()
		uuidStr = "fooo"
		for key, value := range attributes.keyValue {
			log.Println("Adding:", key, " and value ", value, " to " , uuidStr)
			attrMapper.AddQueryToUUID(key,value,uuidStr)
		}
		return uuidStr
	} else {
		log.Fatalf("Could not generate GUID for %v \n Error %v \n", attributes, err)
	}
	return ""
}

func init() {
	AttrMapperManagerInjector = *NewAttrMapperManager()
	AttrMapperManagerInjector.Set("default", NewMemAttrMapper())
}
