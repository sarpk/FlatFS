// In memory attribute mapper

package first

import (
	"log"
	"github.com/nu7hatch/gouuid"
)

type MemAttrMapper struct {
	AttrMapper
	queryToUuid map[string]map[string][]string
}

func NewMemAttrMapper() *MemAttrMapper {
	memAttrMapper := &MemAttrMapper{
		queryToUuid: make(map[string]map[string][]string, 0),
	}
	return memAttrMapper
}

func (attrMapper *MemAttrMapper) AddQueryToUUID(key, value, uuid string) {
	if attrMapper.queryToUuid[key] == nil {
		attrMapper.queryToUuid[key] = make(map[string][]string, 0)
	}
	attrMapper.queryToUuid[key][value] = append(attrMapper.queryToUuid[key][value], uuid)
}


//Known issue: If there are more than one match, then it also assumes it doesn't exist
func (attrMapper *MemAttrMapper) GetAddedUUID(attributes *QueryKeyValue) (string, bool) {
	uniqueVal := make(map[string]bool, 0)
	itemAddedToMap := false
	for key, value := range attributes.keyValue {
		if attrMapper.queryToUuid[key] == nil || attrMapper.queryToUuid[key][value] == nil {
			return "", false
		}
		if (len(uniqueVal) == 0 && !itemAddedToMap) {
			for _, uuid := range attrMapper.queryToUuid[key][value] {
				uniqueVal[uuid] = true
				itemAddedToMap = true
			}
		} else {
			for _, uuid := range attrMapper.queryToUuid[key][value] {
				lessUniqueVals := make(map[string]bool, 0)
				if uniqueVal[uuid] {
					lessUniqueVals[uuid] = true
				}
				uniqueVal = lessUniqueVals
			}
		}
		if(len(uniqueVal) == 0 && itemAddedToMap) { //it must mean that it's not unique enough
			return "", false
		}
	}
	if (len(uniqueVal) == 1) {
		for uniqueUuid := range uniqueVal {
			return uniqueUuid, true
		}
	}
	return "", false
}

func (attrMapper *MemAttrMapper) CreateNewUUID(attributes *QueryKeyValue) string {
	if uuid, err := uuid.NewV4(); err == nil {
		uuidStr := uuid.String()
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

func (attrMapper *MemAttrMapper) CreateFromQuery(attributes *QueryKeyValue) string {
	log.Println("Mocking middleware")
	uuidStr, attributeAdded := attrMapper.GetAddedUUID(attributes)
	if attributeAdded {
		return uuidStr
	}
	return attrMapper.CreateNewUUID(attributes)
}

func init() {
	AttrMapperManagerInjector = *NewAttrMapperManager()
	AttrMapperManagerInjector.Set("default", NewMemAttrMapper())
}
