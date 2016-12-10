// In memory attribute mapper

package FlatFS

import (
	"testing"
	"path"
)

func TestFileCreateAndListWithExactName(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent := "This is the test string"
	attr1 := "foo=hello"
	attr2 := "bar=world"
	exactPath := path.Join(mountPoint, attr1 + "," + attr2)

	write_to_file(exactPath, testContent)
	lsContent := exec_cmd("ls -l " + exactPath)

	assert_string_contains(lsContent, attr1, t)
	assert_string_contains(lsContent, attr2, t)
	Terminate(mountPoint)
}

func TestFileCreateAndListWithFirstAttr(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent := "This is the test string"
	attr1 := "foo=hello"
	attr2 := "bar=world"
	exactPath := path.Join(mountPoint, attr1 + "," + attr2)

	write_to_file(exactPath, testContent)
	lsContent := exec_cmd("ls -l " + path.Join(mountPoint, "?" + attr1))

	assert_string_contains(lsContent, attr1, t)
	assert_string_contains(lsContent, attr2, t)
	Terminate(mountPoint)
}

func TestFileCreateAndListWithSecondAttr(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent := "This is the test string"
	attr1 := "foo=hello"
	attr2 := "bar=world"
	exactPath := path.Join(mountPoint, attr1 + "," + attr2)

	write_to_file(exactPath, testContent)
	lsContent := exec_cmd("ls -l " + path.Join(mountPoint, "?" + attr2))

	assert_string_contains(lsContent, attr1, t)
	assert_string_contains(lsContent, attr2, t)
	Terminate(mountPoint)
}

