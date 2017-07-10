package FlatFS

import (
	"testing"
	"io/ioutil"
	"log"
	"github.com/sarpk/FlatFS/flatFS/test/benchmarking/utils"
	"time"
)

func TestSetup(t *testing.T) {
	UtilsFlatFs.RecurseThroughFolders(UtilsFlatFs.MOUNT_POINT_PATH, t)
}

func TestFileLookUpBenchmarkForFlatFS(t *testing.T) {
	start := time.Now()
	FileLookupBenchmark("FlatFsFileNames.txt")
	log.Printf("TestFileLookUpBenchmarkForFlatFS took %s", time.Since(start))
}

func TestFileLookUpBenchmarkForHFS(t *testing.T) {
	start := time.Now()
	FileLookupBenchmark("HFSFileNames.txt")
	log.Printf("TestFileLookUpBenchmarkForHFS took %s", time.Since(start))
}


func TestTerminate(t *testing.T) {
	UtilsFlatFs.Terminate()
}

func FileLookupBenchmark(fileName string) {
	UtilsFlatFs.FileBenchmark(fileName, ReadFileWrapper)
}

func ReadFileWrapper(fileName string) {
	_, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println("Read file wrapper err ", err)
	}
}
