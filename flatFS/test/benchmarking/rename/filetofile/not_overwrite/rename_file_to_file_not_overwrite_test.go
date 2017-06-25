package FlatFS

import (
	"testing"
	"github.com/sarpk/FlatFS/flatFS/test/benchmarking/utils"
	"log"
	"time"
)


func TestSetupNotOverwrite(t *testing.T) {
	UtilsFlatFs.RecurseThroughFolders("/tmp/lpbckmtpt/", UtilsFlatFs.MOUNT_POINT_PATH, t)
}

func TestFileToFileMoveNotOverwriteBenchmarkForHFS(t *testing.T) {
	start := time.Now()
	UtilsFlatFs.FileBenchmarkTwoProcess("HFSFileNames.txt", RenameFileWrapperAppendWithRandomStr)
	log.Printf("TestFileToFileMoveNotOverwriteBenchmarkForHFS took %s", time.Since(start))
}

func TestFileToFileMoveNotOverwriteBenchmarkForFlatFs(t *testing.T) {
	start := time.Now()
	UtilsFlatFs.FileBenchmarkTwoProcess("FlatFsFileNames.txt", RenameFileWrapperAppendWithRandomStr)
	log.Printf("TestFileToFileMoveNotOverwriteBenchmarkForFlatFs took %s", time.Since(start))
}

func RenameFileWrapperAppendWithRandomStr(oldPath, newPath string) {
	newPath += UtilsFlatFs.RAND_STR
	UtilsFlatFs.RenameFileWrapper(oldPath, newPath)
}
