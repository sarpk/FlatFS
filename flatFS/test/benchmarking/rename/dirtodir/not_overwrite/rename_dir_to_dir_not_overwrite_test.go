package FlatFS

import (
	"testing"
	"github.com/sarpk/FlatFS/flatFS/test/benchmarking/utils"
	"log"
	"time"
	"os"
)

func TestSetup(t *testing.T) {
	UtilsFlatFs.RecurseThroughFolderWithSave(UtilsFlatFs.MOUNT_POINT_PATH, t, true)
}

func TestDirToDirectoryMoveForHFS(t *testing.T) {
	start := time.Now()
	UtilsFlatFs.FileBenchmarkTwoProcess("HFSFileNames.txt", RenameFileWrapperAddRandomPath)
	log.Printf("TestDirToDirectoryMoveForHFS took %s", time.Since(start))
}

func TestDirToDirectoryMoveForFlatFs(t *testing.T) {
	start := time.Now()
	UtilsFlatFs.FileBenchmarkTwoProcess("FlatFsFileNames.txt", RenameFileWrapperAppendWithQueryRandom)
	log.Printf("TestDirToDirectoryMoveForFlatFs took %s", time.Since(start))
}

func TestTerminate(t *testing.T) {
	UtilsFlatFs.Terminate()
}

func RenameFileWrapperAppendWithQueryRandom(oldPath, newPath string) {
	UtilsFlatFs.RenameFileWrapper(UtilsFlatFs.AppendQueryParam(oldPath), UtilsFlatFs.AppendQueryParam(newPath) + ",random:" + UtilsFlatFs.RAND_STR)
}

func RenameFileWrapperAddRandomPath(oldPath, newPath string) {
	newRandomPath := newPath + UtilsFlatFs.RAND_STR + "/"
	os.Mkdir(newRandomPath, os.ModePerm)
	UtilsFlatFs.RenameFileWrapper(oldPath, newRandomPath)
}
