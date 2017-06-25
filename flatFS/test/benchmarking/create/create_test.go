package FlatFS

import (
	"testing"
	"log"
	"os"
	"github.com/sarpk/FlatFS/flatFS/test/benchmarking/utils"
	"time"
)

func TestSetup(t *testing.T) {
	UtilsFlatFs.RecurseThroughFolders("/tmp/lpbckmtpt/", UtilsFlatFs.MOUNT_POINT_PATH, t)
}

func TestFileCreateBenchmarkForFlatFS(t *testing.T) {
	start := time.Now()
	FileCreateBenchmark("FlatFsFileNames.txt")
	log.Printf("TestFileCreateBenchmarkForFlatFS took %s", time.Since(start))
}

func TestFileCreateBenchmarkForHFS(t *testing.T) {
	start := time.Now()
	FileCreateBenchmark("HFSFileNames.txt")
	log.Printf("TestFileCreateBenchmarkForHFS took %s", time.Since(start))
}

func FileCreateBenchmark(fileName string) {
	UtilsFlatFs.FileBenchmark(fileName, CreateFileWrapper)
}

func CreateFileWrapper(fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Println("Create file wrapper err ", err)
	}
	defer file.Close()
}