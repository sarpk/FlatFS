// A Go mirror of libfuse's hello.c

package main

import (
	"flag"
	"log"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"os"
	"path/filepath"
	"syscall"
	"io"
)

type HelloFs struct {
	pathfs.FileSystem
}

//func (me *HelloFs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
//	log.Printf("getattr name is %s", name)
//	switch name {
//	case "file.txt":
//		return &fuse.Attr{
//			Mode: fuse.S_IFREG | 0644, Size: uint64(len(name)),
//		}, fuse.OK
//	case "":
//		return &fuse.Attr{
//			Mode: fuse.S_IFDIR | 0755,
//		}, fuse.OK
//	}
//	return nil, fuse.ENOENT
//}

func (me *HelloFs) GetAttr(name string, context *fuse.Context) (a *fuse.Attr, code fuse.Status) {
	fullPath := me.GetPath(name)
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

	f, err := os.Open(me.GetPath(name))
	if err != nil {
		return nil, fuse.ToStatus(err)
	}
	want := 500
	output := make([]fuse.DirEntry, 0, want)
	for {
		infos, err := f.Readdir(want)
		for i := range infos {
			// workaround forhttps://code.google.com/p/go/issues/detail?id=5960
			if infos[i] == nil {
				continue
			}
			n := infos[i].Name()
			d := fuse.DirEntry{
				Name: n,
			}
			if s := fuse.ToStatT(infos[i]); s != nil {
				d.Mode = uint32(s.Mode)
			} else {
				log.Printf("ReadDir entry %q for %q has no stat info", n, name)
			}
			output = append(output, d)
		}
		if len(infos) < want || err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Readdir() returned err:", err)
			break
		}
	}
	f.Close()

	return output, fuse.OK
}

func (me *HelloFs) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	log.Printf("open namea is %s", name)
	f, err := os.OpenFile(me.GetPath(name), int(flags), 0)
	if err != nil {
		return nil, fuse.ToStatus(err)
	}
	return nodefs.NewLoopbackFile(f), fuse.OK

	//
	//if name != "file.txt" {
	//	return nil, fuse.ENOENT
	//}
	//if flags & fuse.O_ANYWRITE != 0 {
	//	return nil, fuse.EPERM
	//}
	//return nodefs.NewDataFile([]byte("asdasd")), fuse.OK
}

func (me *HelloFs) GetPath(relPath string) string {
	return filepath.Join("/tmp/firstDir", relPath)
}

func (me *HelloFs) Create(name string, flags uint32, mode uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	log.Printf("create file name is %s", name)

	f, err := os.OpenFile(me.GetPath(name), int(flags)|os.O_CREATE, os.FileMode(mode))
	return nodefs.NewLoopbackFile(f), fuse.ToStatus(err)
}

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatal("Usage:\n  hello MOUNTPOINT")
	}
	nfs := pathfs.NewPathNodeFs(&HelloFs{FileSystem: pathfs.NewDefaultFileSystem()}, nil)
	nfs.SetDebug(true)
	server, _, err := nodefs.MountRoot(flag.Arg(0), nfs.Root(), nil)
	server.SetDebug(true)
	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}
	server.Serve()
}
