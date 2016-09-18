// In memory attribute mapper

package first

import (
	"log"
)

type MemAttrMapper struct {
	AttrMapper
	queryToUuid          map[string]map[string][]string
	uuidToAttributeValue map[string]map[string]string
}

func NewMemAttrMapper() *MemAttrMapper {
	memAttrMapper := &MemAttrMapper{
		queryToUuid: make(map[string]map[string][]string, 0),
		uuidToAttributeValue: make(map[string]map[string]string, 0),
	}
	return memAttrMapper
}

func (attrMapper *MemAttrMapper) AddQueryToUUID(key, value, uuid string) {
	if attrMapper.queryToUuid[key] == nil {
		attrMapper.queryToUuid[key] = make(map[string][]string, 0)
	}
	attrMapper.queryToUuid[key][value] = append(attrMapper.queryToUuid[key][value], uuid)

	if attrMapper.uuidToAttributeValue[uuid] == nil {
		attrMapper.uuidToAttributeValue[uuid] = make(map[string]string, 0)
	}
	attrMapper.uuidToAttributeValue[uuid][key] = value
}

func IsQueryDoesntExistInTheAttributeMap(strings map[string]map[string][]string, key string, value string) bool {
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

func AttributesEqual(uniqueResAttrs map[string]string, attributes map[string]string) bool {
	if len(uniqueResAttrs) != len(attributes) {
		return false
	}
	for attr := range uniqueResAttrs {
		if _, ok := attributes[attr]; !ok {
			return false
		}
	}
	return true
}

func (attrMapper *MemAttrMapper) ReturnEqualAttributeResult(uniqueVal map[string]bool, attributes map[string]string) (string, bool) {
	if uniqueVal == nil {
		return "", false
	}
	for uniqueUuid := range uniqueVal {
		if AttributesEqual(attrMapper.uuidToAttributeValue[uniqueUuid], attributes) {
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
			return "", false
		}
		uniqueVal = ReduceUniqueValueMapFromAttributeMapper(attrMapper.queryToUuid[key][value], uniqueVal)
		if len(uniqueVal) == 0 {
			//it must mean that it's not unique enough
			return "", false
		}
	}
	if len(uniqueVal) == 0 {
		return "", false //No unique UUID for the given query found
	}

	path, unique := attrMapper.ReturnEqualAttributeResult(uniqueVal, attributes.keyValue)
	if unique {
		return path, true
	}
	return "", true //It's not unique so it's a folder
}

func (attrMapper *MemAttrMapper) FindAllMatchingMultipleUUIDs(attributes *QueryKeyValue) (map[string]bool, bool) {
	uniqueVal, found := attrMapper.ReturnFirstUUIDFromAttribute(attributes.keyValue)
	if !found {
		return nil, false
	}
	for key, value := range attributes.keyValue {
		if IsQueryDoesntExistInTheAttributeMap(attrMapper.queryToUuid, key, value) {
			return nil, false
		}
		uniqueVal = ReduceUniqueValueMapFromAttributeMapper(attrMapper.queryToUuid[key][value], uniqueVal)
		if len(uniqueVal) == 0 {
			//it must mean that it's not unique enough
			return nil, false
		}
	}
	if len(uniqueVal) > 0 {
		return uniqueVal, true
	}
	return nil, false
}

func (attrMapper *MemAttrMapper) FindAllMatchingQueries(attributes *QueryKeyValue) ([]UUIDToQuery, bool) {
	uuids, found := attrMapper.FindAllMatchingMultipleUUIDs(attributes)
	if !found {
		return nil, false
	}
	queryKeyValues := []UUIDToQuery{}
	for uuid := range uuids {
		queryKeyValue := QueryKeyValue{
			attrMapper.uuidToAttributeValue[uuid],
		}

		foundQuery := UUIDToQuery{
			uuid,
			queryKeyValue,
		}
		queryKeyValues = append(queryKeyValues, foundQuery)
	}
	return queryKeyValues, true
}

func (attrMapper *MemAttrMapper) Close() {
	//TODO Save it to disk
}

func (attrMapper *MemAttrMapper) CreateFromQuery(attributes *QueryKeyValue) string {
	log.Println("Mocking middleware")
	uuidStr, attributeAdded := attrMapper.GetAddedUUID(attributes)
	if attributeAdded {
		return uuidStr
	}
	return CreateNewUUID(attributes, attrMapper.AddQueryToUUID)
}
