package FlatFS

import (
	"flag"
	"log"

	"github.com/sarpk/go-fuse/fuse"
	"github.com/sarpk/go-fuse/fuse/nodefs"
	"github.com/sarpk/go-fuse/fuse/pathfs"
	"os"
	"path/filepath"
	"syscall"
	"strings"
	"fmt"
	"bytes"
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

func (me *FlatFs) GetFileInfoFromUUID(uuid string) os.FileInfo {
	file, oErr := os.Open(me.GetPath(uuid))
	if oErr != nil {
		return nil
	}
	fInfo, sErr := file.Stat()
	if sErr != nil {
		return nil
	}
	return fInfo
}

func ConvertToString(query QueryKeyValue) string {
	var result bytes.Buffer
	for key, value := range query.keyValue {
		if result.Len() != 0 {
			result.WriteString(",")

		}
		result.WriteString(fmt.Sprintf("%v=%v", key, value))
	}
	return result.String()
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

func (me *FlatFs) GetPath(relPath string) string {
	return filepath.Join(me.flatStorage, relPath)
}

func (me *FlatFs) Unlink(name string, context *fuse.Context) (code fuse.Status) {
	parsedQuery, isFile := ParseQuery(name)
	if !isFile {
		return fuse.EINVAL;
	}
	return me.UnlinkParsedQuery(parsedQuery)
}

func (me *FlatFs) UnlinkParsedQuery(parsedQuery *QueryKeyValue) fuse.Status {
	uuid, fileFound := me.attrMapper.GetAddedUUID(parsedQuery, true)
	if !fileFound {
		return fuse.ENODATA;
	}
	fullPath := me.GetPath(uuid)
	deleteStatus := fuse.ToStatus(syscall.Unlink(fullPath))
	if deleteStatus == fuse.OK {
		me.attrMapper.DeleteUUIDFromQuery(parsedQuery, uuid)
	}
	return deleteStatus
}

func (me *FlatFs) Mkdir(path string, mode uint32, context *fuse.Context) (code fuse.Status) {
	return fuse.ToStatus(os.Mkdir(me.GetPath(path), os.FileMode(mode)))
}

type UUIDToQuery struct {
	uuid string
	querykeyValue QueryKeyValue
}

type QueryKeyValue struct {
	keyValue map[string]string
}

func NewQueryKeyValue() *QueryKeyValue {
	return &QueryKeyValue{
		keyValue: make(map[string]string, 0),
	}
}

func ParseQuery(raw string) (*QueryKeyValue, bool) {
	isFile := true
	if strings.IndexByte(raw, '?') == 0  {
		isFile = false
		raw = raw[1:]
	}
	query := NewQueryKeyValue()
	for _, kv := range strings.Split(raw, ",") {
		pair := strings.Split(kv, "=")
		if (len(pair) == 2) {
			query.keyValue[pair[0]] = pair[1]
		}
	}
	return query, isFile
}

func (me *FlatFs) CreateWithNewPath(name string, flags uint32, mode uint32, context *fuse.Context) (file nodefs.File, code fuse.Status, newPath string) {
	log.Printf("create file name is %s", name)
	parsedQuery, _ := ParseQuery(name)
	newPath = me.attrMapper.CreateFromQuery(parsedQuery)
	log.Printf("Saving the file name as %s", newPath)
	f, err := os.OpenFile(me.GetPath(newPath), int(flags) | os.O_CREATE, os.FileMode(mode))
	return nodefs.NewLoopbackFile(f), fuse.ToStatus(err), newPath
}

var (
	AttrMapperManagerInjector AttrMapperManager
)

func Prepare() {
	AttrMapperManagerInjector = *NewAttrMapperManager()
	AttrMapperManagerInjector.Set("default", NewMemAttrMapper())
	AttrMapperManagerInjector.Set("sqlite", NewSQLiteAttrMapper())
}

func Start() {
	Prepare()
	flag.Parse()
	if len(flag.Args()) < 3 {
		log.Fatal("Usage:\n  FlatFS MOUNTPOINT FLATSTORAGE [backend] \n  [backend] can be 'default' (in memory) or 'sqlite' ")
	}
	attrMapperFromManager:= AttrMapperManagerInjector.Get(flag.Arg(2))
	defer attrMapperFromManager.Close()
	flatFs := &FlatFs{
		FileSystem: pathfs.NewDefaultFileSystem(),
		attrMapper: attrMapperFromManager,
		flatStorage: flag.Arg(1),
	}
	nfs := pathfs.NewPathNodeFs(flatFs, nil)
	nfs.SetDebug(true)
	server, _, err := nodefs.MountRoot(flag.Arg(0), nfs.Root(), nil)
	server.SetDebug(true)
	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}
	server.Serve()
}
