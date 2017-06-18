package FlatFS

import (
	"testing"
	"log"
	"os"
	"strings"
)

var RAND_STR = "MN12tA_j"

func TestSetupForFlatFS(t *testing.T) {
	TestFileCreateBenchmarkForFlatFS(t)
	TestFileDeleteBenchmarkForFlatFS(t)
}

func TestSetupForHFS(t *testing.T) {
	TestFileCreateBenchmarkForHFS(t)
	TestFileDeleteBenchmarkForHFS(t)
}

func TestFileToFileMoveOverwriteBenchmarkForHFS(t *testing.T) {
	TestFileCreateBenchmarkForHFS(t)
	defer TestFileDeleteBenchmarkForHFS(t)
	FileBenchmarkTwoProcess("HFSFileNames.txt", RenameFileWrapper)
}

func TestFileToFileMoveOverwriteBenchmarkForFlatFS(t *testing.T) {
	TestFileCreateBenchmarkForFlatFS(t)
	defer TestFileDeleteBenchmarkForFlatFS(t)
	FileBenchmarkTwoProcess("FlatFsFileNames.txt", RenameFileWrapper)
}

func TestFileToFileMoveNotOverwriteBenchmarkForHFS(t *testing.T) {
	TestFileCreateBenchmarkForHFS(t)
	defer TestFileDeleteBenchmarkForHFS(t)
	FileBenchmarkTwoProcess("HFSFileNames.txt", RenameFileWrapperAppendWithRandomStr)
}

func TestFileToFileMoveNotOverwriteBenchmarkForFlatFs(t *testing.T) {
	TestFileCreateBenchmarkForFlatFS(t)
	defer TestFileDeleteBenchmarkForFlatFS(t)
	FileBenchmarkTwoProcess("FlatFsFileNames.txt", RenameFileWrapperAppendWithRandomStr)
}

func TestFileToDirectoryMoveForFlatFs(t *testing.T) {
	TestFileCreateBenchmarkForFlatFS(t)
	defer TestFileDeleteBenchmarkForFlatFS(t)
	FileBenchmarkTwoProcess("FlatFsFileNames.txt", RenameFileWrapperAppendWithQuery)
}

func TestFileToDirectoryMoveForHFS(t *testing.T) {
	TestFileCreateBenchmarkForHFS(t)
	defer TestFileDeleteBenchmarkForHFS(t)
	FileBenchmarkTwoProcess("HFSFileNames.txt", RenameFileWrapperRemoveLastPath)
}

func TestDirToDirectoryMoveForHFS(t *testing.T) {
	TestFileCreateBenchmarkForHFS(t)
	defer TestFileDeleteBenchmarkForHFS(t)
	FileBenchmarkTwoProcess("HFSFileNames.txt", RenameFileWrapperAddRandomPath)
}

func TestDirToDirectoryMoveForFlatFs(t *testing.T) {
	TestFileCreateBenchmarkForFlatFS(t)
	defer TestFileDeleteBenchmarkForFlatFS(t)
	FileBenchmarkTwoProcess("FlatFsFileNames.txt", RenameFileWrapperAppendWithQueryRandom)
}

func TestDirToParentDirMoveForHFS(t *testing.T) {
	TestFileCreateBenchmarkForHFS(t)
	defer TestFileDeleteBenchmarkForHFS(t)
	FileBenchmarkTwoProcess("HFSFileNames.txt", RenameFileWrapperToParent)
}

func TestDirToParentDirMoveForFlatFs(t *testing.T) {
	TestFileCreateBenchmarkForFlatFS(t)
	defer TestFileDeleteBenchmarkForFlatFS(t)
	FileBenchmarkTwoProcess("FlatFsFileNames.txt", RenameFileWrapperRemoveLastAttribute)
}

func RenameFileWrapperRemoveLastAttribute(oldPath, newPath string) {
	RenameFileWrapper(oldPath, DeleteQuery(GetLastAttributeValue(oldPath)))
}

func RenameFileWrapperToParent(oldPath, newPath string) {
	RenameFileWrapper(oldPath, GetParentFolder(oldPath))
}

func RenameFileWrapperAppendWithQueryRandom(oldPath, newPath string) {
	RenameFileWrapper(oldPath, AppendQueryParam(newPath) + ",random:" + RAND_STR)
}

func RenameFileWrapperAddRandomPath(oldPath, newPath string) {
	RenameFileWrapper(oldPath, GetParentFolder(newPath) + RAND_STR + "/")
}

func RenameFileWrapperRemoveLastPath(oldPath, newPath string) {
	RenameFileWrapper(oldPath, GetParentFolder(newPath))
}

func GetLastAttributeValue(path string) (string, string) {
	pairs := strings.Split(path, ",")
	lastPair := pairs[len(pairs) - 1]
	pair := strings.Split(lastPair, ":")
	return pair[0], pair[1]
}

func GetParentFolder(path string) string {
	amountOfDirs := strings.Split(path, "/")
	lastFileLen := len(amountOfDirs[len(amountOfDirs) - 1])
	return path[0:len(path) - lastFileLen]
}

func RenameFileWrapperAppendWithQuery(oldPath, newPath string) {
	RenameFileWrapper(oldPath, AppendQueryParam(newPath))
}

func DeleteQuery(attribute, value string) string {
	return MOUNT_POINT_PATH + "-" + attribute + ":" + value
}

func AppendQueryParam(path string) string {
	return MOUNT_POINT_PATH + "?" + strings.Split(path, MOUNT_POINT_PATH)[1]
}

func RenameFileWrapperAppendWithRandomStr(oldPath, newPath string) {
	newPath += RAND_STR
	RenameFileWrapper(oldPath, newPath)
}

type processFileRename func(string, string)

func FileBenchmarkTwoProcess(fileListPath string, fun processFileRename) {
	fileList := ReadArrays(fileListPath)
	fileListSize := len(fileList);
	j := fileListSize - 1;
	for i := 0; i < fileListSize; i++ {
		if i == j {
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
