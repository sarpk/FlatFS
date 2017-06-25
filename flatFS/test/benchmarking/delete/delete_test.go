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

func TestFileDeleteBenchmarkForFlatFS(t *testing.T) {
	start := time.Now()
	FileDeleteBenchmark("FlatFsFileNames.txt")
	log.Printf("TestFileDeleteBenchmarkForFlatFS took %s", time.Since(start))
}

func TestFileDeleteBenchmarkForHFS(t *testing.T) {
	start := time.Now()
	FileDeleteBenchmark("HFSFileNames.txt")
	log.Printf("TestFileDeleteBenchmarkForHFS took %s", time.Since(start))
}

func FileDeleteBenchmark(fileName string) {
	UtilsFlatFs.FileBenchmark(fileName, DeleteFileWrapper)
}

func DeleteFileWrapper(fileName string) {
	err := os.Remove(fileName)
	if err != nil {
		log.Println("Delete file wrapper err ", err)
	}
}
