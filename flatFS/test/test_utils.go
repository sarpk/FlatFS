// In memory attribute mapper

package FlatFS

import (
	"os"
	"path/filepath"
	"log"
	"fmt"
	"sync"
	"strings"
	"os/exec"
	"github.com/sarpk/FlatFS/flatFS"
	"time"
	"path"
)

func GetCurrentDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("dir is " + dir)
	return dir
}

func CreateFlatStoreDir(dir string) string {
	flatStoreDir := path.Join(dir, "flatDir")
	os.MkdirAll(flatStoreDir, 0777)
	return flatStoreDir
}

func CreateMountPoint(dir string) string {
	mountPointDir := path.Join(dir, "mountpoint")
	os.MkdirAll(mountPointDir, 0777)
	return mountPointDir
}

func CreateFlatFS() string {
	dir := GetCurrentDir()

	mountPointDir := CreateMountPoint(dir)
	flatStoreDir := CreateFlatStoreDir(dir)
	os.Args = append(os.Args, mountPointDir, flatStoreDir, "default")
	go FlatFS.Start() // dispatch in another thread
	time.Sleep(time.Second * 3) // wait 3 secs to start

	return mountPointDir
}

func Terminate(mountPointDir string) {
	exec_cmd("fusermount -u " + mountPointDir, nil)
}

func exec_cmd(cmd string, wg *sync.WaitGroup) {
	fmt.Println("command is ", cmd)
	// splitting head => g++ parts => rest of the command
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("out is %s \n", out)
	//wg.Done() // Need to signal to waitgroup that this goroutine is done
}
