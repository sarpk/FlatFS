// In memory attribute mapper

package FlatFS

import (
	"testing"
	"path"
)

func TestFileCreateAndRead(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent := "This is the test string"
	attr1 := "foo=hello"
	attr2 := "bar=world"
	write_to_file(path.Join(mountPoint, attr1 + "," + attr2), testContent)
	lsContent := exec_cmd("ls -l " + path.Join(mountPoint, "?foo=hello"))
	fileContent := read_from_file(path.Join(mountPoint, "foo=hello,bar=world"))
	assert_string_equals(fileContent, testContent, t)
	assert_string_contains(lsContent, attr1,t)
	assert_string_contains(lsContent, attr2,t)
	Terminate(mountPoint)
}

func TestFileCreateAndListWithExactName(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent := "This is the test string"
	attr1 := "foo=hello"
	attr2 := "bar=world"
	completeAttr := attr1 + "," + attr2
	write_to_file(path.Join(mountPoint, completeAttr), testContent)
	lsContent := exec_cmd("ls -l " + path.Join(mountPoint, completeAttr))
	assert_string_contains(lsContent, attr1,t)
	assert_string_contains(lsContent, attr2,t)
	Terminate(mountPoint)
}

func TestFileCreateAndListWithFirstAttr(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent := "This is the test string"
	attr1 := "foo=hello"
	attr2 := "bar=world"
	write_to_file(path.Join(mountPoint, attr1 + "," + attr2), testContent)
	lsContent := exec_cmd("ls -l " + path.Join(mountPoint, "?" + attr1))
	assert_string_contains(lsContent, attr1,t)
	assert_string_contains(lsContent, attr2,t)
	Terminate(mountPoint)
}

func TestFileCreateAndListWithSecondAttr(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent := "This is the test string"
	attr1 := "foo=hello"
	attr2 := "bar=world"
	write_to_file(path.Join(mountPoint, attr1 + "," + attr2), testContent)
	lsContent := exec_cmd("ls -l " + path.Join(mountPoint, "?" + attr2))
	assert_string_contains(lsContent, attr1,t)
	assert_string_contains(lsContent, attr2,t)
	Terminate(mountPoint)
}

