package FlatFS

import (
	"testing"
	"strings"
	"github.com/sarpk/FlatFS/flatFS/test/benchmarking/utils"
	"log"
	"time"
)

func TestSetup(t *testing.T) {
	UtilsFlatFs.RecurseThroughFolders("/tmp/lpbckmtpt/", UtilsFlatFs.MOUNT_POINT_PATH, t)
}

func TestFileToDirectoryMoveForFlatFs(t *testing.T) {
	start := time.Now()
	UtilsFlatFs.FileBenchmarkTwoProcess("FlatFsFileNames.txt", RenameFileWrapperAppendWithQuery)
	log.Printf("TestFileToDirectoryMoveForFlatFs took %s", time.Since(start))
}

func TestFileToDirectoryMoveForHFS(t *testing.T) {
	start := time.Now()
	UtilsFlatFs.FileBenchmarkTwoProcess("HFSFileNames.txt", RenameFileWrapperRemoveLastPath)
	log.Printf("TestFileToDirectoryMoveForHFS took %s", time.Since(start))
}

func RenameFileWrapperRemoveLastPath(oldPath, newPath string) {
	UtilsFlatFs.RenameFileWrapper(oldPath, GetParentFolder(newPath))
}

func RenameFileWrapperAppendWithQuery(oldPath, newPath string) {
	UtilsFlatFs.RenameFileWrapper(oldPath, AppendQueryParam(newPath))
}

func GetParentFolder(path string) string {
	amountOfDirs := strings.Split(path, "/")
	lastFileLen := len(amountOfDirs[len(amountOfDirs) - 1])
	return path[0:len(path) - lastFileLen] + UtilsFlatFs.RAND_STR
}

func AppendQueryParam(path string) string {
	return UtilsFlatFs.MOUNT_POINT_PATH + "?" + strings.Split(path, UtilsFlatFs.MOUNT_POINT_PATH)[1]
}


