// In memory attribute mapper

package first

import (
	"log"
	"github.com/nu7hatch/gouuid"
)

type AttrMapper interface {
	CreateFromQuery(*QueryKeyValue) string
	GetAddedUUID(attributes *QueryKeyValue) (string, bool)
	FindAllMatchingQueries(attributes *QueryKeyValue) ([]QueryKeyValue, bool)
	Close()
}

func CreateNewUUID(attributes *QueryKeyValue, addQueryToUUID func(string, string, string)) string {
	if uuid, err := uuid.NewV4(); err == nil {
		uuidStr := uuid.String()
		for key, value := range attributes.keyValue {
			log.Println("Adding:", key, " and value ", value, " to ", uuidStr)
			addQueryToUUID(key, value, uuidStr)
		}
		return uuidStr
	} else {
		log.Fatalf("Could not generate GUID for %v \n Error %v \n", attributes, err)
	}
	return ""
}
