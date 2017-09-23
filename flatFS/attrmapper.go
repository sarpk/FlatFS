package FlatFS

import (
	"log"
	"github.com/nu7hatch/gouuid"
	"syscall"
	"github.com/sarpk/go-fuse/fuse"
)

type AttrMapper interface {
	GetAddedUUID(attributes *QueryKeyValue, queryType QueryType) (string, bool)
	FindAllMatchingQueries(attributes *QueryKeyValue) ([]UUIDToQuery, bool)
	DeleteQueryToUUID(key, value, uuid string)
	Close()
	AddQueryToUUID(key, value, uuid string)
}

func RenameQuery(oldSpec *QueryKeyValue, newSpec *QueryKeyValue, fs *FlatFs) {
	uuidMatchingToFile, found := fs.attrMapper.GetAddedUUID(oldSpec, createFileSpecQueryType())
	if !found {
		return
	}
	DeleteUUIDFromQuery(oldSpec, fs.attrMapper.DeleteQueryToUUID, uuidMatchingToFile)
	UnlinkParsedQuery(newSpec, fs)
	AddUUIDToAttributes(newSpec, fs.attrMapper.AddQueryToUUID, uuidMatchingToFile)
}

func AppendOldSpec(oldSpec *QueryKeyValue, newSpec *QueryKeyValue, fs *FlatFs) {
	uuidMatchingToFile, found := fs.attrMapper.GetAddedUUID(oldSpec, createFileSpecQueryType())
	if !found {
		return
	}
	DeleteUUIDFromQuery(oldSpec, fs.attrMapper.DeleteQueryToUUID, uuidMatchingToFile)

	AddUUIDToAttributes(AppendQueryKeyValue(oldSpec, newSpec), fs.attrMapper.AddQueryToUUID, uuidMatchingToFile)
}

func CreateFromQuery(attributes *QueryKeyValue, fs *FlatFs) string {
	uuidStr, attributeAdded := fs.attrMapper.GetAddedUUID(attributes, createFileSpecQueryType())
	if attributeAdded {
		return uuidStr
	}
	return CreateNewUUID(attributes, fs.attrMapper.AddQueryToUUID)
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

func DeleteUUIDFromQuery(attributes *QueryKeyValue, deleteUUIDFromKeyValue func(string, string, string), uuid string) {
	for key, value := range attributes.keyValue {
		deleteUUIDFromKeyValue(key, value, uuid)
	}
}

func AddUUIDToAttributes(attributes *QueryKeyValue, addQueryToUUID func(string, string, string), uuid string) {
	for key, value := range attributes.keyValue {
		addQueryToUUID(key, value, uuid)
	}
}

func UnlinkParsedQuery(parsedQuery *QueryKeyValue, fs *FlatFs) fuse.Status {
	uuid, fileFound := fs.attrMapper.GetAddedUUID(parsedQuery, createFileSpecQueryType())
	if !fileFound {
		return fuse.ENODATA;
	}
	fullPath := fs.GetPath(uuid)
	deleteStatus := fuse.ToStatus(syscall.Unlink(fullPath))
	if deleteStatus == fuse.OK {
		DeleteUUIDFromQuery(parsedQuery, fs.attrMapper.DeleteQueryToUUID, uuid)
	}
	return deleteStatus
}

func AppendQueryKeyValue(toBeAppended *QueryKeyValue, toAppend *QueryKeyValue) *QueryKeyValue {
	for key, value := range toAppend.keyValue {
		toBeAppended.keyValue[key] = value
	}
	return toBeAppended
}
