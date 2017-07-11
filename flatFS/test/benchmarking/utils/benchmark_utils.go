package UtilsFlatFs

import (
	"testing"
	"os"
	"path"
	"io"
	"strings"
	"bytes"
	"fmt"
	"math/rand"
	"encoding/gob"
	"io/ioutil"
	"github.com/sarpk/FlatFS/flatFS/test/in-memory"
	"log"
	"syscall"
)

var FlatFsFileNames = make([]string, 0)
var HFSFileNames = make([]string, 0)

var MOUNT_POINT_PATH = FlatFS.GetCurrentDir() + "/mountpoint/"
var MOUNT_POINT_PATH_HFS = FlatFS.GetCurrentDir() + "/hfs/"

var RAND_STR = "MN12tA_j"

type processFileRename func(string, string)

func RenameFileWrapperAppendWithQueryRandom(oldPath, newPath string) {
	RenameFileWrapper(oldPath, AppendQueryParam(newPath) + ",random:" + RAND_STR)
}

func RenameFileWrapperAddRandomPath(oldPath, newPath string) {
	RenameFileWrapper(oldPath, GetParentFolder(newPath) + RAND_STR + "/")
}

func GetParentFolder(path string) string {
	amountOfDirs := strings.Split(path, "/")
	lastFileLen := len(amountOfDirs[len(amountOfDirs) - 1])
	return path[0:len(path) - lastFileLen]
}

func AppendQueryParam(path string) string {
	splitStr := strings.Split(path, MOUNT_POINT_PATH)
	if len(splitStr) < 2 {
		return ""
	}
	return MOUNT_POINT_PATH + "?" + splitStr[1]
}

func RenameFileWrapper(oldPath, newPath string) {
	if len(oldPath) == 0 || len(newPath) == 0 {
		return
	}
	err := syscall.Rename(oldPath, newPath)
	if err != nil {
		log.Println("Rename file wrapper err ", err, oldPath, newPath)
	}
}

func FileBenchmarkTwoProcess(fileListPath string, fun processFileRename) {
	fileList := ReadArrays(fileListPath)
	fileListSize := len(fileList);
	j := fileListSize - 1;
	for i := 0; i < fileListSize; i++ {
		if i >= j {
			break;
		}
		fun(fileList[i], fileList[j])
		j--
	}
}

type processFile func(string)

func DeleteFileWrapper(fileName string) {
	err := os.Remove(fileName)
	if err != nil {
		log.Println("Delete file wrapper err ", err)
	}
}

func FileBenchmark(fileListPath string, fun processFile) {
	fileList := ReadArrays(fileListPath)
	for _, fileName := range fileList {
		fun(fileName)
	}
}

func RecurseThroughFolders(flatFsPath string, t *testing.T) {
	RecurseThroughFolderWithSave(flatFsPath, t, false)
}

func RecurseThroughFolderWithSave(flatFsPath string, t *testing.T, onlyFolderSave bool) {
	FlatFS.CreateFlatFS()
	CopyFiles()
	rootPath := MOUNT_POINT_PATH_HFS
	rootDirectory := ScanDirectory(rootPath, t)
	nextDirectories := FilterAsDirectoryPath(rootDirectory, rootPath)
	rand.Seed(63)
	SaveFilesInDirectoryToFlatFs(rootDirectory, flatFsPath, rootPath, rootPath, onlyFolderSave)
	if onlyFolderSave {
		FlatFsFileNames = FlatFsFileNames[1:]
	}
	for len(nextDirectories) != 0 {
		nextDirectory := nextDirectories[0]
		nextDirectories = nextDirectories[1:]
		currentScan := ScanDirectory(nextDirectory, t)
		directoriesToAdd := FilterAsDirectoryPath(currentScan, nextDirectory)
		nextDirectories = append(nextDirectories, directoriesToAdd...)
		if onlyFolderSave {
			if len(directoriesToAdd) == 0 {
				SaveFilesInDirectoryToFlatFs(currentScan, flatFsPath, nextDirectory, rootPath, onlyFolderSave)
			}
		} else {
			SaveFilesInDirectoryToFlatFs(currentScan, flatFsPath, nextDirectory, rootPath, onlyFolderSave)
		}

	}
	ShuffleArrays()
	WriteArrays()
	log.Println("size is ", len(FlatFsFileNames))
}

func CopyFiles() {
	os.Mkdir(MOUNT_POINT_PATH_HFS, os.ModePerm)
	err := CopyDir("/tmp/benchmark/", MOUNT_POINT_PATH_HFS)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Directory copied")
	}

}

func Terminate() {
	FlatFS.Terminate(MOUNT_POINT_PATH)
}

func ReadArrays(fileName string) []string {
	b, _ := ioutil.ReadFile(fileName)
	fileList := []string{}
	gob.NewDecoder(bytes.NewBuffer(b)).Decode(&fileList)
	return fileList
}

func WriteArrays() {
	buf1 := &bytes.Buffer{}
	gob.NewEncoder(buf1).Encode(FlatFsFileNames)

	err := ioutil.WriteFile("FlatFsFileNames.txt", buf1.Bytes(), 0644)
	if err != nil {
		panic(err)
	}

	buf2 := &bytes.Buffer{}
	gob.NewEncoder(buf2).Encode(HFSFileNames)
	err = ioutil.WriteFile("HFSFileNames.txt", buf2.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}

func ShuffleArrays() {
	for i := range FlatFsFileNames {
		j := rand.Intn(i + 1)
		FlatFsFileNames[i], FlatFsFileNames[j] = FlatFsFileNames[j], FlatFsFileNames[i]
		HFSFileNames[i], HFSFileNames[j] = HFSFileNames[j], HFSFileNames[i]
	}
}

func SaveFilesInDirectoryToFlatFs(dirContent []os.FileInfo, flatFsPath, currPath, rootPath string, onlyFolder bool) {
	attributesToAdd := currPath[len(rootPath):]
	attributes := strings.Split(attributesToAdd, "/")[1:]
	levelCount := 1
	var attrBuf bytes.Buffer
	attrBuf.WriteString(flatFsPath)
	for _, attribute := range attributes {
		attrBuf.WriteString(fmt.Sprintf("level_%v:%v,", levelCount, attribute))
		levelCount = levelCount + 1
	}

	filesToBeSaved := FilterAsFileNames(dirContent)

	//log.Println("files are going to be saved in: ", attrBuf.String())
	if (onlyFolder) {
		pathToSave := attrBuf.String()
		pathToSave = pathToSave[0:len(pathToSave) - 1]
		FlatFsFileNames = append(FlatFsFileNames, pathToSave)
		HFSFileNames = append(HFSFileNames, currPath)
	}
	for _, fileName := range filesToBeSaved {
		var filePath bytes.Buffer
		filePath.Write(attrBuf.Bytes())
		filePath.WriteString(fmt.Sprintf("level_%v:%v", levelCount, fileName))
		fileNameToSave := filePath.String()
		if rand.Intn(10) == 5 {
			//10% chance to match
			os.Create(fileNameToSave)
			if !onlyFolder {
				FlatFsFileNames = append(FlatFsFileNames, fileNameToSave)
				HFSFileNames = append(HFSFileNames, MOUNT_POINT_PATH_HFS + attributesToAdd + "/" + fileName)
			}
		}
	}
}

func FilterAsDirectoryPath(dirContent []os.FileInfo, currPath string) []string {
	result := make([]string, 0)
	for _, content := range dirContent {
		if content.IsDir() {
			dirPath := path.Join(currPath, content.Name())
			result = append(result, dirPath)
		}
	}
	return result
}

func FilterAsFileNames(dirContent []os.FileInfo) []string {
	result := make([]string, 0)
	for _, content := range dirContent {
		if !content.IsDir() {
			result = append(result, content.Name())
		}
	}
	return result
}

func ScanDirectory(dir string, t *testing.T) []os.FileInfo {

	file, err := os.Open(dir)
	if err != nil {
		t.Fatalf("open %q failed: %v", dir, err)
	}
	defer file.Close()
	files, err2 := file.Readdir(100000)
	if err2 != nil && err2 != io.EOF {
		t.Fatalf("Reading dir %q failed: %v", file, err2)
	}
	return files
}

func CopyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}

	}

	return
}

func CopyDir(source string, dest string) (err error) {

	// get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// create dest dir

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourcefilepointer := source + "/" + obj.Name()

		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			// create sub-directories - recursively
			err = CopyDir(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// perform copy
			err = CopyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
	return
}















