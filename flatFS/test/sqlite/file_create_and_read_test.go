package FlatFS

import (
	"testing"
	"path"
)

func TestFileCreateAndRead(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent := "This is the test string"
	attr1 := "foo:hello"
	attr2 := "bar:world"
	exactPath := path.Join(mountPoint, attr1 + "," + attr2)

	write_to_file(exactPath, testContent)
	lsContent := exec_cmd("ls -l " + path.Join(mountPoint, "?" + attr1))
	fileContent := read_from_file(exactPath)

	assert_string_equals(fileContent, testContent, t)
	assert_string_contains(lsContent, attr1, t)
	assert_string_contains(lsContent, attr2, t)
	Terminate(mountPoint)
}

func TestFileCreateAndReadWithThreeAttrs(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent := "This is the test string"
	attr1 := "foo:hello"
	attr2 := "bar:world"
	attr3 := "flat:fs"
	exactPath := path.Join(mountPoint, attr1 + "," + attr2 + "," + attr3)

	write_to_file(exactPath, testContent)
	lsContent := exec_cmd("ls -l " + path.Join(mountPoint, "?" + attr1))
	fileContent := read_from_file(exactPath)

	assert_string_equals(fileContent, testContent, t)
	assert_string_contains(lsContent, attr1, t)
	assert_string_contains(lsContent, attr2, t)
	assert_string_contains(lsContent, attr3, t)
	Terminate(mountPoint)
}

func TestFileCreateAndReadWithThreeAttrsShuffleOrder(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent := "This is the test string"
	attr1 := "foo:hello"
	attr2 := "bar:world"
	attr3 := "flat:fs"
	exactPath := path.Join(mountPoint, attr2 + "," + attr3 + "," + attr1)

	write_to_file(exactPath, testContent)
	lsContent := exec_cmd("ls -l " + path.Join(mountPoint, "?" + attr3))
	fileContent := read_from_file(path.Join(mountPoint, attr3 + "," + attr1 + "," + attr2))

	assert_string_equals(fileContent, testContent, t)
	assert_string_contains(lsContent, attr1, t)
	assert_string_contains(lsContent, attr2, t)
	assert_string_contains(lsContent, attr3, t)
	Terminate(mountPoint)
}

func TestMultipleFileCreateAndRead(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent1 := "This is the test string"
	testContent2 := "This is another file"

	attr1 := "foo:hello"
	attr2 := "bar:world"
	attr3 := "flat:fs"
	exactPath1 := path.Join(mountPoint, attr1 + "," + attr2)
	exactPath2 := path.Join(mountPoint, attr2 + "," + attr3)

	write_to_file(exactPath1, testContent1)
	lsContent := exec_cmd("ls -l " + path.Join(mountPoint, "?" + attr2))
	fileContent := read_from_file(exactPath1)

	assert_string_equals(fileContent, testContent1, t)
	assert_string_contains(lsContent, attr1, t)
	assert_string_contains(lsContent, attr2, t)
	assert_string_not_contains(lsContent, attr3, t)

	write_to_file(exactPath2, testContent2)
	lsContent = exec_cmd("ls -l " + path.Join(mountPoint, "?" + attr2))
	fileContent = read_from_file(exactPath1)

	assert_string_equals(fileContent, testContent1, t) //Check that it's not overwritten
	fileContent = read_from_file(exactPath2)
	assert_string_equals(fileContent, testContent2, t)

	assert_string_contains(lsContent, attr1, t)
	assert_string_contains(lsContent, attr2, t)
	assert_string_contains(lsContent, attr3, t)

	Terminate(mountPoint)
}

func TestMultipleFileCreateAndReadWithInclusiveAttributesLessOrder(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent1 := "This is the test string"
	testContent2 := "This is another file"

	attr1 := "foo:hello"
	attr2 := "bar:world"
	attr3 := "flat:fs"
	exactPath1 := path.Join(mountPoint, attr1 + "," + attr2)
	exactPath2 := path.Join(mountPoint, attr1 + "," + attr2 + "," + attr3)

	write_to_file(exactPath1, testContent1)
	lsContent := exec_cmd("ls -l " + path.Join(mountPoint, "?" + attr2))
	fileContent := read_from_file(exactPath1)

	assert_string_equals(fileContent, testContent1, t)
	assert_string_contains(lsContent, attr1, t)
	assert_string_contains(lsContent, attr2, t)
	assert_string_not_contains(lsContent, attr3, t)

	write_to_file(exactPath2, testContent2)
	lsContent = exec_cmd("ls -l " + path.Join(mountPoint, "?" + attr2))
	fileContent = read_from_file(exactPath1)

	assert_string_equals(fileContent, testContent1, t) //Check that it's not overwritten
	fileContent = read_from_file(exactPath2)
	assert_string_equals(fileContent, testContent2, t)

	assert_string_contains(lsContent, attr1, t)
	assert_string_contains(lsContent, attr2, t)
	assert_string_contains(lsContent, attr3, t)

	Terminate(mountPoint)
}


func TestMultipleFileCreateAndReadWithInclusiveAttributesHigherOrder(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent1 := "This is the test string"
	testContent2 := "This is another file"

	attr1 := "foo:hello"
	attr2 := "bar:world"
	attr3 := "flat:fs"
	exactPath1 := path.Join(mountPoint, attr1 + "," + attr2 + "," + attr3)
	exactPath2 := path.Join(mountPoint, attr1 + "," + attr2)

	write_to_file(exactPath1, testContent1)
	lsContent := exec_cmd("ls -l " + path.Join(mountPoint, "?" + attr2))
	fileContent := read_from_file(exactPath1)

	assert_string_equals(fileContent, testContent1, t)
	assert_string_contains(lsContent, attr1, t)
	assert_string_contains(lsContent, attr2, t)
	assert_string_contains(lsContent, attr3, t)

	write_to_file(exactPath2, testContent2)
	lsContent = exec_cmd("ls -l " + path.Join(mountPoint, "?" + attr2))
	fileContent = read_from_file(exactPath1)

	assert_string_equals(fileContent, testContent1, t) //Check that it's not overwritten
	fileContent = read_from_file(exactPath2)
	assert_string_equals(fileContent, testContent2, t)

	assert_string_contains(lsContent, attr1, t)
	assert_string_contains(lsContent, attr2, t)
	assert_string_contains(lsContent, attr3, t)

	Terminate(mountPoint)
}
