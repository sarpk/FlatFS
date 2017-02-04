package FlatFS

import (
	"log"

	"github.com/sarpk/go-fuse/fuse"
	"github.com/sarpk/go-fuse/fuse/nodefs"
	"github.com/sarpk/go-fuse/fuse/pathfs"
	"os"
	"syscall"
	"github.com/nu7hatch/gouuid"
)

type FlatFs struct {
	pathfs.FileSystem
	attrMapper  AttrMapper
	flatStorage string
}

func (flatFs *FlatFs) GetAttrWithPath(path, fileName string, context *fuse.Context) (*fuse.Attr, fuse.Status, string) {
	log.Printf("GetAttrWithPath for path is %s and name is %s", path, fileName)
	pathQuery, pathQueryType := ParseQuery(path)
	fileNameQuery, fileNameQueryType := ParseQuery(fileName)
	if (QueryOrAdd(pathQueryType) || QueryOrAdd(fileNameQueryType)) {
		for k, v := range pathQuery.keyValue {
			fileNameQuery.keyValue[k] = v
		}
		log.Printf("PathQuery keyvalue is %v", pathQuery.keyValue)
		log.Printf("fileNameQuery keyvalue is %v", fileNameQuery.keyValue)
	}
	fullPath := ConvertToString(*fileNameQuery)
	if (pathQueryType.fileSpec || fileNameQueryType.fileSpec) {
		attr, code := flatFs.GetAttr(fullPath, context)
		return attr, code, fullPath
	} else if QueryOrAdd(pathQueryType) {
		attr, code := flatFs.GetAttr(path, context)
		return attr, code, fullPath
	} else if QueryOrAdd(fileNameQueryType) {
		attr, code := flatFs.GetAttr(fileName, context)
		return attr, code, fullPath
	}

	attr, code := flatFs.AttrStatHandle(fullPath, fullPath, fileNameQueryType)
	return attr, code, fullPath
}

func (flatFs *FlatFs) GetAttr(name string, context *fuse.Context) (a *fuse.Attr, code fuse.Status) {
	log.Printf("GetAttr for name is %s", name)
	//fullPath :=
	parsedQueryPath, queryType := ParseQuery(name)
	fullPath, fileFound := flatFs.attrMapper.GetAddedUUID(parsedQueryPath, queryType)
	if !fileFound {
		fullPath = name
	}
	fullPath = flatFs.GetPath(fullPath)
	log.Printf("Found path is  %s", fullPath)
	return flatFs.AttrStatHandle(name, fullPath, queryType)

}

func (FlatFs *FlatFs) AttrStatHandle(originalPath, fullPath string, queryType QueryType) (*fuse.Attr, fuse.Status) {
	var err error = nil
	st := syscall.Stat_t{}
	if originalPath == "" || !queryType.fileSpec {
		//also check if type is file
		// When GetAttr is called for the toplevel directory, we always want
		// to look through symlinks.
		err = syscall.Stat(fullPath, &st)
	} else {
		err = syscall.Lstat(fullPath, &st)
	}

	if err != nil {
		return nil, fuse.ToStatus(err)
	}
	a := &fuse.Attr{}
	a.FromStat(&st)
	return a, fuse.OK
}

func (flatFs *FlatFs) RenameWithNewPath(oldName string, newPath, newName string, context *fuse.Context) fuse.Status {
	log.Printf("RenameWithNewPath %s to path %s and name %s", oldName, newPath, newName)
	newPathQuery, newPathqueryType := ParseQuery(newPath)
	if QueryOrAdd(newPathqueryType) {
		newNameQuery, _ := ParseQuery(newName)
		for k, v := range newPathQuery.keyValue {
			newNameQuery.keyValue[k] = v
		}
		newName = ConvertToString(*newNameQuery)
	}
	return flatFs.Rename(oldName, newName, context)
}

func (flatFs *FlatFs) Rename(oldName string, newName string, context *fuse.Context) (code fuse.Status) {
	log.Printf("Renaming %s to %s", oldName, newName)
	oldSpec, _ := ParseQuery(oldName)
	newSpec, queryType := ParseQuery(newName)
	if !queryType.fileSpec {
		flatFs.attrMapper.AppendOldSpec(oldSpec, newSpec, flatFs)
	} else {
		flatFs.attrMapper.RenameQuery(oldSpec, newSpec, flatFs)
	}
	return fuse.OK
}

func (flatFs *FlatFs) OpenDir(name string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	log.Printf("opendira name is %s", name)
	parsedQuery, _ := ParseQuery(name)
	foundQueries, fileFound := flatFs.attrMapper.FindAllMatchingQueries(parsedQuery)
	if !fileFound {
		_, err := os.Open(flatFs.GetPath(name))
		return nil, fuse.ToStatus(err)
	}

	for something, another := range foundQueries {
		log.Printf("will process this query %v and %v", something, another)
	}
	output := make([]fuse.DirEntry, 0)
	for _, foundQuery := range foundQueries {
		d := fuse.DirEntry{
			Name: ConvertToString(foundQuery.querykeyValue),
		}
		if s := fuse.ToStatT(flatFs.GetFileInfoFromUUID(foundQuery.uuid)); s != nil {
			d.Mode = uint32(s.Mode)
		}
		output = append(output, d)
	}

	return output, fuse.OK
}

func (flatFs *FlatFs) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	log.Printf("open namea is %s", name)
	_, parseError := uuid.ParseHex(name)
	if parseError == nil {
		//It's a valid UUID so just parse that
		return flatFs.OpenFileAsLoopback(name, int(flags))
	}
	name, _ = flatFs.attrMapper.GetAddedUUID(ParseQuery(name))
	return flatFs.OpenFileAsLoopback(name, int(flags))
}

func (flatFs *FlatFs) Unlink(name string, context *fuse.Context) (code fuse.Status) {
	parsedQuery, querySpec := ParseQuery(name)
	if !querySpec.fileSpec {
		return fuse.EINVAL;
	}
	return flatFs.UnlinkParsedQuery(parsedQuery)
}

func (flatFs *FlatFs) Mkdir(path string, mode uint32, context *fuse.Context) (code fuse.Status) {
	return fuse.ToStatus(os.Mkdir(flatFs.GetPath(path), os.FileMode(mode)))
}

func (flatFs *FlatFs) CreateWithNewPath(name string, flags uint32, mode uint32, context *fuse.Context) (file nodefs.File, code fuse.Status, newPath string) {
	log.Printf("create file name is %s", name)
	parsedQuery, _ := ParseQuery(name)
	newPath = flatFs.attrMapper.CreateFromQuery(parsedQuery)
	log.Printf("Saving the file name as %s", newPath)
	f, err := os.OpenFile(flatFs.GetPath(newPath), int(flags) | os.O_CREATE, os.FileMode(mode))
	return nodefs.NewLoopbackFile(f), fuse.ToStatus(err), newPath
}


