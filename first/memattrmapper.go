// In memory attribute mapper

package first

import (
	"log"
	"github.com/nu7hatch/gouuid"
)

type MemAttrMapper struct {
	AttrMapper
	queryToUuid         map[string]map[string][]string
	uuidToAttributeName map[string][]string
}

func NewMemAttrMapper() *MemAttrMapper {
	memAttrMapper := &MemAttrMapper{
		queryToUuid: make(map[string]map[string][]string, 0),
		uuidToAttributeName: make(map[string][]string, 0),
	}
	return memAttrMapper
}

func (attrMapper *MemAttrMapper) AddQueryToUUID(key, value, uuid string) {
	if attrMapper.queryToUuid[key] == nil {
		attrMapper.queryToUuid[key] = make(map[string][]string, 0)
	}
	attrMapper.queryToUuid[key][value] = append(attrMapper.queryToUuid[key][value], uuid)
	attrMapper.uuidToAttributeName[uuid] = append(attrMapper.uuidToAttributeName[uuid], key)
}

func IsQueryDoesntExistInTheAttributeMap(strings map[string]map[string][]string, key string, value string) {
	return strings == nil || strings[key] == nil || strings[key][value] == nil
}

func (attrMapper *MemAttrMapper) ReturnFirstUUIDFromAttribute(strings map[string]string) (map[string]bool, bool) {
	uniqueVal := make(map[string]bool, 0)
	for key, value := range strings {
		if IsQueryDoesntExistInTheAttributeMap(attrMapper.queryToUuid, key, value) {
			return nil, false
		}
		for _, uuid := range attrMapper.queryToUuid[key][value] {
			uniqueVal[uuid] = true
		}
		return uniqueVal, true
	}
	return nil, false
}

func ReduceUniqueValueMapFromAttributeMapper(queryToUuid []string, uniqueVal map[string]bool) map[string]bool {
	lessUniqueVals := make(map[string]bool, 0)
	for _, uuid := range queryToUuid {
		if uniqueVal[uuid] {
			lessUniqueVals[uuid] = true
		}
	}
	return lessUniqueVals
}

func AttributesEqual(uniqueResAttrs []string, attributes map[string]string) bool {
	if len(uniqueResAttrs) != len(attributes) {
		return false
	}
	for attr := range uniqueResAttrs {
		if attributes[attr] == nil {
			return false
		}
	}
	return true
}

func (attrMapper *MemAttrMapper) ReturnEqualAttributeResult(uniqueVal map[string]bool, attributes map[string]string) (string, bool) {
	for uniqueUuid := range uniqueVal {
		if AttributesEqual(attrMapper.uuidToAttributeName[uniqueUuid], attributes) {
			return uniqueUuid, true
		}
	}
	return "", false
}

func ReturnFirstForMap(uniqueVal map[string]bool) (string, bool) {
	for uniqueUuid := range uniqueVal {
		return uniqueUuid, true
	}
	return "", false
}

func (attrMapper *MemAttrMapper) GetAddedUUID(attributes *QueryKeyValue) (string, bool) {
	uniqueVal, found := attrMapper.ReturnFirstUUIDFromAttribute(attributes.keyValue)
	if !found {
		return "", false
	}
	for key, value := range attributes.keyValue {
		if IsQueryDoesntExistInTheAttributeMap(attrMapper.queryToUuid, key, value) {
			return nil, false
		}
		uniqueVal = ReduceUniqueValueMapFromAttributeMapper(attrMapper.queryToUuid[key][value], uniqueVal)
		if len(uniqueVal) == 0 {
			//it must mean that it's not unique enough
			return "", false
		}
	}
	if len(uniqueVal) > 1 {
		return attrMapper.ReturnEqualAttributeResult(uniqueVal, attributes.keyValue)
	}
	return ReturnFirstForMap(uniqueVal)
}

func (attrMapper *MemAttrMapper) CreateNewUUID(attributes *QueryKeyValue) string {
	if uuid, err := uuid.NewV4(); err == nil {
		uuidStr := uuid.String()
		for key, value := range attributes.keyValue {
			log.Println("Adding:", key, " and value ", value, " to ", uuidStr)
			attrMapper.AddQueryToUUID(key, value, uuidStr)
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
