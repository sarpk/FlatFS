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

func TestFileToFileMoveOverwriteBenchmarkForHFS(t *testing.T) {
	FileBenchmarkTwoProcess("HFSFileNames.txt", RenameFileWrapper)
}

func TestFileToFileMoveOverwriteBenchmarkForFlatFS(t *testing.T) {
	FileBenchmarkTwoProcess("FlatFsFileNames.txt", RenameFileWrapper)
}

type processFileRename func(string, string)


func FileBenchmarkTwoProcess(fileListPath string, fun processFileRename) {
	fileList := ReadArrays(fileListPath)
	fileListSize := len(fileList);
	j := fileListSize-1;
	for i := 0; i < fileListSize; i++ {
		if i==j {
			break;
		}
		fun(fileList[i], fileList[j])
		j--
	}
}

func RenameFileWrapper(oldPath, newPath string) {
	err := os.Rename(oldPath, newPath)
	if err != nil {
		log.Println("Rename file wrapper err ", err)
	}
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

func FileBenchmark(fileListPath string, fun processFile) {
	fileList := ReadArrays(fileListPath)
	for _, fileName := range fileList {
		fun(fileName)
	}
}


func CreateFileWrapper(fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Println("Create file wrapper err ", err)
	}
	defer file.Close()
}

func DeleteFileWrapper(fileName string) {
	err := os.Remove(fileName)
	if err != nil {
		log.Println("Delete file wrapper err ", err)
	}
}

func ReadFileWrapper(fileName string) {
	_, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println("Read file wrapper err ", err)
	}
}
