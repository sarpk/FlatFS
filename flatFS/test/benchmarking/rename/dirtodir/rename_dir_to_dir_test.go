package FlatFS

import (
	"testing"
	"github.com/sarpk/FlatFS/flatFS/test/benchmarking/utils"
	"log"
	"time"
	"os"
)

func TestSetup(t *testing.T) {
	UtilsFlatFs.RecurseThroughFolders("/tmp/lpbckmtpt/", UtilsFlatFs.MOUNT_POINT_PATH, t)
}

func TestDirToDirectoryMoveForHFS(t *testing.T) {
	UtilsFlatFs.FileBenchmarkTwoProcess("HFSFileNames.txt", RenameFileWrapperAddRandomPath)
}

func TestDirToDirectoryMoveForFlatFs(t *testing.T) {
	start := time.Now()
	UtilsFlatFs.FileBenchmarkTwoProcess("FlatFsFileNames.txt", RenameFileWrapperAppendWithQueryRandom)
	log.Printf("TestDirToDirectoryMoveForFlatFs took %s", time.Since(start))
}

func RenameFileWrapperAppendWithQueryRandom(oldPath, newPath string) {
	UtilsFlatFs.RenameFileWrapper(oldPath, UtilsFlatFs.AppendQueryParam(newPath) + ",random:" + UtilsFlatFs.RAND_STR)
}

func RenameFileWrapperAddRandomPath(oldPath, newPath string) {
	newRandomPath := UtilsFlatFs.GetParentFolder(newPath) + UtilsFlatFs.RAND_STR + "/"
	os.Mkdir(newRandomPath, os.ModePerm)
	UtilsFlatFs.RenameFileWrapper(UtilsFlatFs.GetParentFolder(oldPath), newRandomPath)
}
