package FlatFS

import (
	"testing"
	"os"
	"path"
	"io"
	"strings"
	"bytes"
	"fmt"
	"math/rand"
)

var FlatFsFileNames = make([]string, 0)
var HFSFileNames = make([]string, 0)


func RecurseThroughFolders(rootPath, flatFsPath string, t *testing.T) {
	rootDirectory := ScanDirectory(rootPath, t)
	nextDirectories := FilterAsDirectoryPath(rootDirectory, rootPath)
	rand.Seed(63)
	SaveFilesInDirectoryToFlatFs(rootDirectory, flatFsPath, rootPath, rootPath)
	for len(nextDirectories) != 0 {
		nextDirectory := nextDirectories[0]
		nextDirectories = nextDirectories[1:]
		currentScan := ScanDirectory(nextDirectory, t)
		directoriesToAdd := FilterAsDirectoryPath(currentScan, nextDirectory)
		nextDirectories = append(nextDirectories, directoriesToAdd...)
		SaveFilesInDirectoryToFlatFs(currentScan, flatFsPath, nextDirectory, rootPath)
	}
	ShuffleArrays()
	//log.Println("size is ", len(FlatFsFileNames))
}

func ShuffleArrays() {
	for i := range FlatFsFileNames {
		j := rand.Intn(i + 1)
		FlatFsFileNames[i], FlatFsFileNames[j] = FlatFsFileNames[j], FlatFsFileNames[i]
		HFSFileNames[i], HFSFileNames[j] = HFSFileNames[j], HFSFileNames[i]
	}
}

func SaveFilesInDirectoryToFlatFs(dirContent []os.FileInfo, flatFsPath, currPath, rootPath string) {
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

	for _, fileName := range filesToBeSaved {
		var filePath bytes.Buffer
		filePath.Write(attrBuf.Bytes())
		filePath.WriteString(fmt.Sprintf("level_%v:%v", levelCount, fileName))
		fileNameToSave := filePath.String()
		//os.Create(fileNameToSave)
		if rand.Intn(10) == 5 { //10% chance to match
			FlatFsFileNames = append(FlatFsFileNames, fileNameToSave)
			HFSFileNames = append(HFSFileNames, attributesToAdd+ "/" +fileNameToSave)
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


















