// In memory attribute mapper

package FlatFS

import (
	"testing"
	"path"
)

func TestFileCreateAndListWithExactName(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent := "This is the test string"
	attr1 := "foo:hello"
	attr2 := "bar:world"
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
	attr1 := "foo:hello"
	attr2 := "bar:world"
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
	attr1 := "foo:hello"
	attr2 := "bar:world"
	exactPath := path.Join(mountPoint, attr1 + "," + attr2)

	write_to_file(exactPath, testContent)
	lsContent := exec_cmd("ls -l " + path.Join(mountPoint, "?" + attr2))

	assert_string_contains(lsContent, attr1, t)
	assert_string_contains(lsContent, attr2, t)
	Terminate(mountPoint)
}

func TestMultipleFileCreateAndListWithSingularFileQuery(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent := "Test Content"

	attr1 := "foo:hello"
	attr2 := "bar:world"
	attr3 := "flat:fs"
	exactPath1 := path.Join(mountPoint, attr1 + "," + attr2)
	exactPath2 := path.Join(mountPoint, attr1 + "," + attr2 + "," + attr3)
	exactPath3 := path.Join(mountPoint, attr3)

	write_to_file(exactPath1, testContent)
	write_to_file(exactPath2, testContent)
	write_to_file(exactPath3, testContent)

	lsContent := exec_cmd("ls -l " + path.Join(mountPoint, attr1))
	assert_string_contains(lsContent, "", t)
	assert_string_not_contains(lsContent, attr1, t)
	assert_string_not_contains(lsContent, attr2, t)
	assert_string_not_contains(lsContent, attr3, t)

	lsContent = exec_cmd("ls -l " + path.Join(mountPoint, attr2))
	assert_string_contains(lsContent, "", t)
	assert_string_not_contains(lsContent, attr1, t)
	assert_string_not_contains(lsContent, attr2, t)
	assert_string_not_contains(lsContent, attr3, t)

	lsContent = exec_cmd("ls -l " + path.Join(mountPoint, attr3 + "," + attr2))
	assert_string_contains(lsContent, "", t)
	assert_string_not_contains(lsContent, attr1, t)
	assert_string_not_contains(lsContent, attr2, t)
	assert_string_not_contains(lsContent, attr3, t)

	lsContent = exec_cmd("ls -l " + exactPath1)
	assert_string_contains(lsContent, attr1, t)
	assert_string_contains(lsContent, attr2, t)
	assert_string_not_contains(lsContent, attr3, t)

	lsContent = exec_cmd("ls -l " + exactPath2)
	assert_string_contains(lsContent, attr1, t)
	assert_string_contains(lsContent, attr2, t)
	assert_string_contains(lsContent, attr3, t)

	lsContent = exec_cmd("ls -l " + exactPath3)
	assert_string_not_contains(lsContent, attr1, t)
	assert_string_not_contains(lsContent, attr2, t)
	assert_string_contains(lsContent, attr3, t)

	Terminate(mountPoint)
}

func TestMultipleFileCreateAndListWithMultiFileQuery(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent := "Test Content"

	attr1 := "foo:hello"
	attr2 := "bar:world"
	attr3 := "flat:fs"

	listDelim := "?"

	exactPath1 := path.Join(mountPoint, attr1 + "," + attr2)
	exactPath2 := path.Join(mountPoint, attr1 + "," + attr2 + "," + attr3)
	exactPath3 := path.Join(mountPoint, attr3)

	write_to_file(exactPath1, testContent)
	write_to_file(exactPath2, testContent)
	write_to_file(exactPath3, testContent)

	lsContent := exec_cmd("ls -l " + path.Join(mountPoint, listDelim + attr1))
	lsContent = assert_string_contains_per_line(lsContent, []string{attr1, attr2, attr3}, t)
	assert_string_contains_per_line(lsContent, []string{attr1, attr2}, t)

	lsContent = exec_cmd("ls -l " + path.Join(mountPoint, listDelim + attr2))
	lsContent = assert_string_contains_per_line(lsContent, []string{attr1, attr2, attr3}, t)
	assert_string_contains_per_line(lsContent, []string{attr1, attr2}, t)

	lsContent = exec_cmd("ls -l " + path.Join(mountPoint, listDelim + attr3 + "," + attr2))
	assert_string_contains_per_line(lsContent, []string{attr1, attr2, attr3}, t)

	lsContent = exec_cmd("ls -l " + path.Join(mountPoint, listDelim + attr1 + "," + attr2))
	lsContent = assert_string_contains_per_line(lsContent, []string{attr1, attr2, attr3}, t)
	assert_string_contains_per_line(lsContent, []string{attr1, attr2}, t)

	lsContent = exec_cmd("ls -l " + path.Join(mountPoint, listDelim + attr3 + "," + attr2 + "," + attr1))
	assert_string_contains_per_line(lsContent, []string{attr1, attr2, attr3}, t)

	lsContent = exec_cmd("ls -l " + path.Join(mountPoint, listDelim + attr3))
	lsContent = assert_string_contains_per_line(lsContent, []string{attr1, attr2, attr3}, t)
	assert_string_contains_per_line(lsContent, []string{attr3}, t)

	Terminate(mountPoint)
}
