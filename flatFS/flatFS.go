package FlatFS

import (
	"log"

	"github.com/sarpk/go-fuse/fuse"
	"github.com/sarpk/go-fuse/fuse/nodefs"
	"github.com/sarpk/go-fuse/fuse/pathfs"
	"os"
	"syscall"
	"io/ioutil"
)

type FlatFs struct {
	pathfs.FileSystem
	attrMapper AttrMapper
	flatStorage string
}

func (flatFs *FlatFs) GetAttr(name string, context *fuse.Context) (a *fuse.Attr, code fuse.Status) {
	log.Printf("GetAttr for name is %s", name)
	//fullPath :=
	fullPath, fileFound := flatFs.attrMapper.GetAddedUUID(ParseQuery(name))
	if !fileFound {
		fullPath = name
	}
	fullPath = flatFs.GetPath(fullPath)
	log.Printf("Found path is  %s", fullPath)
	var err error = nil
	st := syscall.Stat_t{}
	if name == "" {
		// When GetAttr is called for the toplevel directory, we always want
		// to look through symlinks.
		err = syscall.Stat(fullPath, &st)
	} else {
		err = syscall.Lstat(fullPath, &st)
	}

	if err != nil {
		return nil, fuse.ToStatus(err)
	}
	a = &fuse.Attr{}
	a.FromStat(&st)
	return a, fuse.OK
}

func (flatFs *FlatFs) Rename(oldName string, newName string, context *fuse.Context) (code fuse.Status) {
	log.Printf("Renaming %s to %s", oldName, newName)
	oldSpec, _ := ParseQuery(oldName)
	newSpec, isNewSpecAFile := ParseQuery(newName)
	if !isNewSpecAFile {
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
	name, _ = flatFs.attrMapper.GetAddedUUID(ParseQuery(name))
	filePath := flatFs.GetPath(name)
	f, err := os.OpenFile(filePath, int(flags), 0)
	if err != nil {
		return nil, fuse.ToStatus(err)
	}
	syncErr := f.Sync()
	if syncErr != nil {
		return nil, fuse.ToStatus(err)
	}
	//nodefs.NewDataFile(f)
	rdByte, err := ioutil.ReadFile(filePath)
	file = nodefs.NewLoopbackFile(f)
	return nodefs.NewDataFile(rdByte), fuse.OK
	//file, code = flatFs.FileSystem.Open(name,flags,context)
	//if file != nil {
	//	return nil, fuse.ToStatus(err)
	//}
	//return
	//return nodefs.NewLoopbackFile(f), fuse.OK
}

func (flatFs *FlatFs) Unlink(name string, context *fuse.Context) (code fuse.Status) {
	parsedQuery, isFile := ParseQuery(name)
	if !isFile {
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


