// A Go mirror of libfuse's hello.c

package first

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



type DBMiddleware interface {
	FileAttributes(string) string
}

type HelloFs struct {
	pathfs.FileSystem
	attrMapper AttrMapper
}

type MockMiddleware struct {
	DBMiddleware
}

func NewMockMiddleware() *MockMiddleware {
	mockMiddleware := &MockMiddleware{}
	return mockMiddleware
}

func (mockMiddleware *MockMiddleware) FileAttributes(attributes string) string {
	fileAttributes := attributes
	log.Println("Mocking middleware")
	return fileAttributes
}

type DBMiddlewareManager struct {
	middlewares map[string]DBMiddleware
}

func NewDBMiddlewareManager() *DBMiddlewareManager {
	dbMiddlewareManager := &DBMiddlewareManager{
		middlewares: make(map[string]DBMiddleware, 0),
	}
	return dbMiddlewareManager
}

func (dbMiddlewareManager *DBMiddlewareManager) Map() map[string]DBMiddleware {
	return dbMiddlewareManager.middlewares
}

func (dbMiddlewareManager *DBMiddlewareManager) Has(id string) bool {
	_, ok := dbMiddlewareManager.middlewares[id]
	return ok
}

func (dbMiddlewareManager *DBMiddlewareManager) Get(id string) DBMiddleware {
	if dbMiddleware, ok := dbMiddlewareManager.middlewares[id]; ok {
		return dbMiddleware
	}
	log.Fatalf("Implementation %v not found!\n", id)
	return nil
}

func (dbMiddlewareManager *DBMiddlewareManager) Set(id string, dbMiddleware DBMiddleware) DBMiddleware {
	dbMiddlewareManager.middlewares[id] = dbMiddleware
	return dbMiddlewareManager.middlewares[id]
}

func (dbMiddlewareManager *DBMiddlewareManager) Pop(id string) DBMiddleware {
	tempDbMiddleware, ok := dbMiddlewareManager.middlewares[id]
	if ok {
		delete(dbMiddlewareManager.middlewares, id)
	}
	return tempDbMiddleware
}

func (me *HelloFs) GetAttr(name string, context *fuse.Context) (a *fuse.Attr, code fuse.Status) {
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

func (me *HelloFs) OpenDir(name string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	log.Printf("opendira name is %s", name)
	//if name == "" {
	//	c = []fuse.DirEntry{{Name: "file.txt", Mode: fuse.S_IFREG}}
	//	return c, fuse.OK
	//}
	//
	//return nil, fuse.ENOENT
	foundQueries, fileFound := me.attrMapper.FindAllMatchingQueries(ParseQuery(name))
	if !fileFound {
		_, err := os.Open(me.GetPath(name))
		return nil, fuse.ToStatus(err)
	}

	for something, another := range foundQueries {
		log.Printf("will process this query %v and %v", something, another)
	}
	//f, err := os.Open(me.GetPath(name))
	//if err != nil {
	//	return nil, fuse.ToStatus(err)
	//}
	//want := 500
	output := make([]fuse.DirEntry, 0)
	//for {
	//	infos, err := f.Readdir(want)
	//	for i := range infos {
	//		// workaround forhttps://code.google.com/p/go/issues/detail?id=5960
	//		if infos[i] == nil {
	//			continue
	//		}
	//		n := infos[i].Name()
	//		d := fuse.DirEntry{
	//			Name: n,
	//		}
	//		if s := fuse.ToStatT(infos[i]); s != nil {
	//			d.Mode = uint32(s.Mode)
	//		} else {
	//			log.Printf("ReadDir entry %q for %q has no stat info", n, name)
	//		}
	//		output = append(output, d)
	//	}
	//	if len(infos) < want || err == io.EOF {
	//		break
	//	}
	//	if err != nil {
	//		log.Println("Readdir() returned err:", err)
	//		break
	//	}
	//}
	//f.Close()

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

func (me *HelloFs) GetFileInfoFromUUID(uuid string) os.FileInfo {
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

func (me *HelloFs) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	log.Printf("open namea is %s", name)
	name, _ = me.attrMapper.GetAddedUUID(ParseQuery(name))
	f, err := os.OpenFile(me.GetPath(name), int(flags), 0)
	if err != nil {
		return nil, fuse.ToStatus(err)
	}
	return nodefs.NewLoopbackFile(f), fuse.OK

}

func (me *HelloFs) GetPath(relPath string) string {
	return filepath.Join("/tmp/firstDir", relPath)
}

func (me *HelloFs) Unlink(name string, context *fuse.Context) (code fuse.Status) {
	return fuse.ToStatus(syscall.Unlink(me.GetPath(name)))
}

func (me *HelloFs) Mkdir(path string, mode uint32, context *fuse.Context) (code fuse.Status) {
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
	if strings.IndexByte(raw, '/') == 0  {
		isFile = false
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

func (me *HelloFs) CreateWithNewPath(name string, flags uint32, mode uint32, context *fuse.Context) (file nodefs.File, code fuse.Status, newPath string) {
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
	testFunc()
	testFunc()
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatal("Usage:\n  hello MOUNTPOINT")
	}
	attrMapperFromManager:= AttrMapperManagerInjector.Get("default")
	defer attrMapperFromManager.Close()
	helloFs := &HelloFs{
		FileSystem: pathfs.NewDefaultFileSystem(),
		attrMapper: attrMapperFromManager,
	}
	nfs := pathfs.NewPathNodeFs(helloFs, nil)
	nfs.SetDebug(true)
	server, _, err := nodefs.MountRoot(flag.Arg(0), nfs.Root(), nil)
	server.SetDebug(true)
	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}
	server.Serve()
}
