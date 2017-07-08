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

func TestCreateHalfFolder(t *testing.T) {
	start := time.Now()
	UtilsFlatFs.FileBenchmarkTwoProcess("HFSFileNames.txt", CreateNewPathAgain)
	log.Printf("TestCreateHalfFolder took %s", time.Since(start))
}

func TestDirToDirectoryMoveForHFS(t *testing.T) {
	start := time.Now()
	UtilsFlatFs.FileBenchmarkTwoProcess("HFSFileNames.txt", RenameFileWrapperAddOnEmptyFolder)
	log.Printf("TestDirToDirectoryMoveForHFS took %s", time.Since(start))
}

func TestDirToDirectoryMoveForFlatFs(t *testing.T) {
	start := time.Now()
	UtilsFlatFs.FileBenchmarkTwoProcess("FlatFsFileNames.txt", RenameFileWrapperAppend)
	log.Printf("TestDirToDirectoryMoveForFlatFs took %s", time.Since(start))
}

func TestTerminate(t *testing.T) {
	UtilsFlatFs.Terminate()
}

func RenameFileWrapperAppend(oldPath, newPath string) {
	UtilsFlatFs.RenameFileWrapper(UtilsFlatFs.AppendQueryParam(oldPath), UtilsFlatFs.AppendQueryParam(newPath))
}

func CreateNewPathAgain(oldPath, newPath string) {
	log.Printf("Deleting path %s", newPath)
	os.RemoveAll(newPath)
	os.Mkdir(newPath, os.ModePerm)
}

func RenameFileWrapperAddOnEmptyFolder(oldPath, newPath string) {
	UtilsFlatFs.RenameFileWrapper(oldPath, newPath)
}
