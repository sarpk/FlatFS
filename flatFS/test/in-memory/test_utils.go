package FlatFS

import (
	"os"
	"path/filepath"
	"log"
	"fmt"
	"strings"
	"os/exec"
	"github.com/sarpk/FlatFS/flatFS"
	"time"
	"path"
	"io/ioutil"
	"testing"
	"bytes"
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
	time.Sleep(time.Second * 2) // wait 2 secs to finish
	exec_cmd("fusermount -u " + mountPointDir)
	os.Remove("file_metadata.db")
}

func exec_cmd(cmd string) string {
	fmt.Println("command is ", cmd)
	// splitting head => g++ parts => rest of the command
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]
	fmt.Println("parts are  ", parts)
	out, err := exec.Command(head, parts...).Output()
	if err != nil {
		fmt.Printf("err is %s \n", err)
	}
	fmt.Printf("out is %s \n", out)
	return string(out)
}

func write_to_file(filePath, text string) {
	fmt.Println("Creating this file ", filePath)
	f, err := os.OpenFile(filePath, os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		panic(err)
	}
}

func read_from_file(filePath string) string {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		//Do something
	}
	return string(content)
}

func assert_string_equals(str1, str2 string, t *testing.T) {
	if !strings.EqualFold(str1, str2) {
		fmt.Printf("Failing because %s and %s not equal \n", str1, str2)
		t.Fail()
	}
}

func assert_string_contains_per_line(multiLine string, arr[] string, t *testing.T) string {
	fail := true
	var result bytes.Buffer
	for _, str1 := range strings.Split(multiLine, "\n") {
		numToDecrease := len(arr)
		for _, str2 := range arr {
			if strings.Contains(str1, str2) {
				numToDecrease--
			}
		}
		if numToDecrease == 0 {
			fail = false
		} else {
			result.WriteString(str1)
		}
	}

	if fail {
		fmt.Printf("Failing because %s is not included in %s  \n", multiLine, arr)
		t.Fail()
	}
	return result.String()
}

func assert_string_contains(str1, str2 string, t *testing.T) {
	if !strings.Contains(str1, str2) {
		fmt.Printf("Failing because %s is not included in %s  \n", str1, str2)
		t.Fail()
	}
}

func assert_string_not_contains(str1, str2 string, t *testing.T) {
	if strings.Contains(str1, str2) {
		fmt.Printf("Failing because %s is included in %s  \n", str1, str2)
		t.Fail()
	}
}
