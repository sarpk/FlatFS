// In memory attribute mapper

package FlatFS

import (
	"log"
	"github.com/nu7hatch/gouuid"
)

type AttrMapper interface {
	CreateFromQuery(*QueryKeyValue) string
	GetAddedUUID(attributes *QueryKeyValue, queryType QueryType) (string, bool)
	FindAllMatchingQueries(attributes *QueryKeyValue) ([]UUIDToQuery, bool)
	RenameQuery(oldSpec *QueryKeyValue, newSpec *QueryKeyValue, fs *FlatFs)
	AppendOldSpec(oldSpec *QueryKeyValue, newSpec *QueryKeyValue, fs *FlatFs)
	DeleteUUIDFromQuery(attributes *QueryKeyValue, uuid string)
	Close()
}

func AddUUIDToAttributes(attributes *QueryKeyValue, addQueryToUUID func(string, string, string), uuid string) {
	for key, value := range attributes.keyValue {
		log.Println("Adding:", key, " and value ", value, " to ", uuid)
		addQueryToUUID(key, value, uuid)
	}
}

func CreateNewUUID(attributes *QueryKeyValue, addQueryToUUID func(string, string, string)) string {
	if uuid, err := uuid.NewV4(); err == nil {
		uuidStr := uuid.String()
		AddUUIDToAttributes(attributes, addQueryToUUID, uuidStr);
		return uuidStr
	} else {
		log.Fatalf("Could not generate GUID for %v \n Error %v \n", attributes, err)
	}
	return ""
}
