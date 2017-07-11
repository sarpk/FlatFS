package FlatFS

import (
	"testing"
	"log"
	"github.com/sarpk/FlatFS/flatFS/test/benchmarking/utils"
	"time"
)

func TestSetup(t *testing.T) {
	UtilsFlatFs.RecurseThroughFolders(UtilsFlatFs.MOUNT_POINT_PATH, t)
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

func TestTerminate(t *testing.T) {
	UtilsFlatFs.Terminate()
}

func FileDeleteBenchmark(fileName string) {
	UtilsFlatFs.FileBenchmark(fileName, UtilsFlatFs.DeleteFileWrapper)
}


