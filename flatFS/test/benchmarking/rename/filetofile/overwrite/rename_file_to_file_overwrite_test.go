package FlatFS

import (
	"testing"
	"github.com/sarpk/FlatFS/flatFS/test/benchmarking/utils"
	"log"
	"time"
)

func TestSetupOverwrite(t *testing.T) {
	UtilsFlatFs.RecurseThroughFolders(UtilsFlatFs.MOUNT_POINT_PATH, t)
}

func TestFileToFileMoveOverwriteBenchmarkForHFS(t *testing.T) {
	start := time.Now()
	UtilsFlatFs.FileBenchmarkTwoProcess("HFSFileNames.txt", UtilsFlatFs.RenameFileWrapper)
	log.Printf("TestFileToFileMoveOverwriteBenchmarkForHFS took %s", time.Since(start))
}

func TestFileToFileMoveOverwriteBenchmarkForFlatFS(t *testing.T) {
	start := time.Now()
	UtilsFlatFs.FileBenchmarkTwoProcess("FlatFsFileNames.txt", UtilsFlatFs.RenameFileWrapper)
	log.Printf("TestFileToFileMoveOverwriteBenchmarkForFlatFS took %s", time.Since(start))
}

func TestTerminate(t *testing.T) {
	UtilsFlatFs.Terminate()
}


