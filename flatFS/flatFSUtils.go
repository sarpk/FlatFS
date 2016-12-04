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
		log.Print("Given flags are: ", flag.Args())
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

func (flatFs *FlatFs) UnlinkParsedQuery(parsedQuery *QueryKeyValue) fuse.Status {
	uuid, fileFound := flatFs.attrMapper.GetAddedUUID(parsedQuery, true)
	if !fileFound {
		return fuse.ENODATA;
	}
	fullPath := flatFs.GetPath(uuid)
	deleteStatus := fuse.ToStatus(syscall.Unlink(fullPath))
	if deleteStatus == fuse.OK {
		flatFs.attrMapper.DeleteUUIDFromQuery(parsedQuery, uuid)
	}
	return deleteStatus
}

func (flatFs *FlatFs) GetPath(relPath string) string {
	return filepath.Join(flatFs.flatStorage, relPath)
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


func (flatFs *FlatFs) GetFileInfoFromUUID(uuid string) os.FileInfo {
	file, oErr := os.Open(flatFs.GetPath(uuid))
	if oErr != nil {
		return nil
	}
	fInfo, sErr := file.Stat()
	if sErr != nil {
		return nil
	}
	return fInfo
}
