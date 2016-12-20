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

func TestOverwriteFileRename(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent1 := "Test Content"
	testContent2 := "Another File Content"

	attr1 := "foo=hello"
	attr2 := "bar=world"
	attr3 := "flat=fs"

	exactPath1 := path.Join(mountPoint, attr1 + "," + attr2)
	exactPath2 := path.Join(mountPoint, attr2 + "," + attr3)

	write_to_file(exactPath1, testContent1)
	write_to_file(exactPath2, testContent2)
	time.Sleep(time.Second * 1) // TODO fix this wait period!
	exec_cmd("mv " + exactPath1 + " " + exactPath2)
	time.Sleep(time.Second * 1) // TODO fix this wait period!
	fileContent := read_from_file(exactPath2)
	assert_string_equals(fileContent, testContent1, t)

	lsContent := exec_cmd("ls -l " + exactPath1)
	assert_string_not_contains(lsContent, attr1, t)
	assert_string_not_contains(lsContent, attr2, t)
	assert_string_not_contains(lsContent, attr3, t)

	lsContent = exec_cmd("ls -l " + exactPath2)
	lsContent = assert_string_contains_per_line(lsContent, []string{attr2, attr3}, t)
	assert_string_not_contains(lsContent, attr1, t)
	assert_string_not_contains(lsContent, attr2, t)
	assert_string_not_contains(lsContent, attr3, t)

	Terminate(mountPoint)
}
