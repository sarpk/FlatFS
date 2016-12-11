// In memory attribute mapper

package FlatFS

import (
	"testing"
	"path"
)

func TestMultipleFileCreateAndRemove(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent := "Test Content"

	attr1 := "foo=hello"
	attr2 := "bar=world"
	attr3 := "flat=fs"

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

	lsContent = exec_cmd("ls -l " + path.Join(mountPoint, listDelim + attr3))
	lsContent = assert_string_contains_per_line(lsContent, []string{attr1, attr2, attr3}, t)
	assert_string_contains_per_line(lsContent, []string{attr3}, t)

	fileContent := read_from_file(exactPath1)
	assert_string_equals(fileContent, testContent, t)
	exec_cmd("rm " + exactPath1)
	lsContent = exec_cmd("ls -l " + path.Join(mountPoint, listDelim + attr1))
	lsContent = assert_string_contains_per_line(lsContent, []string{attr1, attr2, attr3}, t)
	assert_string_not_contains(lsContent, attr1, t)
	assert_string_not_contains(lsContent, attr2, t)
	fileContent = read_from_file(exactPath1)
	assert_string_equals(fileContent, "", t)

	fileContent = read_from_file(exactPath2)
	assert_string_equals(fileContent, testContent, t)
	exec_cmd("rm " + exactPath2)
	lsContent = exec_cmd("ls -l " + path.Join(mountPoint, listDelim + attr1))
	assert_string_not_contains(lsContent, attr1, t)
	assert_string_not_contains(lsContent, attr2, t)
	assert_string_not_contains(lsContent, attr3, t)
	fileContent = read_from_file(exactPath2)
	assert_string_equals(fileContent, "", t)

	lsContent = exec_cmd("ls -l " + path.Join(mountPoint, listDelim + attr3))
	lsContent = assert_string_contains_per_line(lsContent, []string{attr3}, t)
	assert_string_not_contains(lsContent, attr1, t)
	assert_string_not_contains(lsContent, attr2, t)
	assert_string_not_contains(lsContent, attr3, t)

	fileContent = read_from_file(exactPath3)
	assert_string_equals(fileContent, testContent, t)
	exec_cmd("rm " + exactPath3)
	lsContent = exec_cmd("ls -l " + path.Join(mountPoint, listDelim + attr1))
	assert_string_not_contains(lsContent, attr1, t)
	assert_string_not_contains(lsContent, attr2, t)
	assert_string_not_contains(lsContent, attr3, t)
	fileContent = read_from_file(exactPath3)
	assert_string_equals(fileContent, "", t)

	Terminate(mountPoint)
}
