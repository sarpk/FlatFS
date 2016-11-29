package FlatFS

import (
	"log"

	"github.com/sarpk/go-fuse/fuse"
	"github.com/sarpk/go-fuse/fuse/nodefs"
	"github.com/sarpk/go-fuse/fuse/pathfs"
	"os"
	"syscall"
)

type FlatFs struct {
	pathfs.FileSystem
	attrMapper AttrMapper
	flatStorage string
}

func (me *FlatFs) GetAttr(name string, context *fuse.Context) (a *fuse.Attr, code fuse.Status) {
	log.Printf("GetAttr for name is %s", name)
	//fullPath :=
	fullPath, fileFound := me.attrMapper.GetAddedUUID(ParseQuery(name))
	if !fileFound {
		fullPath = name
	}
	fullPath = me.GetPath(fullPath)
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

func (me *FlatFs) Rename(oldName string, newName string, context *fuse.Context) (code fuse.Status) {
	log.Printf("Renaming %s to %s", oldName, newName)
	oldSpec, _ := ParseQuery(oldName)
	newSpec, isNewSpecAFile := ParseQuery(newName)
	if !isNewSpecAFile {
		me.attrMapper.AppendOldSpec(oldSpec, newSpec, me)
	} else {
		me.attrMapper.RenameQuery(oldSpec, newSpec, me)
	}
	return fuse.OK
}

func (me *FlatFs) OpenDir(name string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	log.Printf("opendira name is %s", name)
	parsedQuery, _ := ParseQuery(name)
	foundQueries, fileFound := me.attrMapper.FindAllMatchingQueries(parsedQuery)
	if !fileFound {
		_, err := os.Open(me.GetPath(name))
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
		if s := fuse.ToStatT(me.GetFileInfoFromUUID(foundQuery.uuid)); s != nil {
			d.Mode = uint32(s.Mode)
		}
		output = append(output, d)
	}

	return output, fuse.OK
}


func (me *FlatFs) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	log.Printf("open namea is %s", name)
	name, _ = me.attrMapper.GetAddedUUID(ParseQuery(name))
	f, err := os.OpenFile(me.GetPath(name), int(flags), 0)
	if err != nil {
		return nil, fuse.ToStatus(err)
	}
	return nodefs.NewLoopbackFile(f), fuse.OK

}

func (me *FlatFs) Unlink(name string, context *fuse.Context) (code fuse.Status) {
	parsedQuery, isFile := ParseQuery(name)
	if !isFile {
		return fuse.EINVAL;
	}
	return me.UnlinkParsedQuery(parsedQuery)
}

func (me *FlatFs) Mkdir(path string, mode uint32, context *fuse.Context) (code fuse.Status) {
	return fuse.ToStatus(os.Mkdir(me.GetPath(path), os.FileMode(mode)))
}

func (me *FlatFs) CreateWithNewPath(name string, flags uint32, mode uint32, context *fuse.Context) (file nodefs.File, code fuse.Status, newPath string) {
	log.Printf("create file name is %s", name)
	parsedQuery, _ := ParseQuery(name)
	newPath = me.attrMapper.CreateFromQuery(parsedQuery)
	log.Printf("Saving the file name as %s", newPath)
	f, err := os.OpenFile(me.GetPath(newPath), int(flags) | os.O_CREATE, os.FileMode(mode))
	return nodefs.NewLoopbackFile(f), fuse.ToStatus(err), newPath
}


