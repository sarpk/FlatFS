// In memory attribute mapper

package FlatFS

import (
	"testing"
	"path"
	"time"
)

func TestSingleFileRename(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent := "Test Content"

	attr1 := "foo=hello"
	attr2 := "bar=world"
	attr3 := "flat=fs"

	initialFilePath := path.Join(mountPoint, attr1 + "," + attr2)
	movedFilePath := path.Join(mountPoint, attr3 + "," + attr2)

	write_to_file(initialFilePath, testContent)
	time.Sleep(time.Second * 1) // TODO fix this wait period!
	exec_cmd("mv " + initialFilePath + " " + movedFilePath)

	fileContent := read_from_file(movedFilePath)
	assert_string_equals(fileContent, testContent, t)

	lsContent := exec_cmd("ls -l " + initialFilePath)
	assert_string_not_contains(lsContent, attr1, t)
	assert_string_not_contains(lsContent, attr2, t)
	assert_string_not_contains(lsContent, attr3, t)

	lsContent = exec_cmd("ls -l " + movedFilePath)
	lsContent = assert_string_contains_per_line(lsContent, []string{attr2, attr3}, t)
	assert_string_not_contains(lsContent, attr1, t)
	assert_string_not_contains(lsContent, attr2, t)
	assert_string_not_contains(lsContent, attr3, t)

	Terminate(mountPoint)
}
