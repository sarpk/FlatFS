package FlatFS

import (
	"testing"
	"io/ioutil"
	"log"
	"os"
)

func TestUtilsFunctions(t *testing.T) {
	RecurseThroughFolders("/tmp/lpbckmtpt/", "/tmp/mountpoint/", t)
}

func TestFileLookUpBenchmarkForFlatFS(t *testing.T) {
	FileLookupBenchmark("FlatFsFileNames.txt")
}

func TestFileLookUpBenchmarkForHFS(t *testing.T) {
	FileLookupBenchmark("HFSFileNames.txt")
}

func TestFileDeleteBenchmarkForFlatFS(t *testing.T) {
	FileDeleteBenchmark("FlatFsFileNames.txt")
}

func TestFileDeleteBenchmarkForHFS(t *testing.T) {
	FileDeleteBenchmark("HFSFileNames.txt")
}

func TestFileCreateBenchmarkForFlatFS(t *testing.T) {
	FileCreateBenchmark("FlatFsFileNames.txt")
}

func TestFileCreateBenchmarkForHFS(t *testing.T) {
	FileCreateBenchmark("HFSFileNames.txt")
}

type processFile func(string)

func FileCreateBenchmark(fileName string) {
	FileBenchmark(fileName, CreateFileWrapper)
}

func FileDeleteBenchmark(fileName string) {
	FileBenchmark(fileName, DeleteFileWrapper)
}

func FileLookupBenchmark(fileName string) {
	FileBenchmark(fileName, ReadFileWrapper)
}

func FileBenchmark(fileName string, fun processFile) {
	fileList := ReadArrays(fileName)
	for _, fileName := range fileList {
		fun(fileName)
	}
}

func CreateFileWrapper(fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
}

func DeleteFileWrapper(fileName string) {
	err := os.Remove(fileName)
	if err != nil {
		log.Println(err)
	}
}

func ReadFileWrapper(fileName string) {
	_, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println(err)
	}
}
