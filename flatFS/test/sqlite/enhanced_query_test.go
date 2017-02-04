// In memory attribute mapper

package FlatFS

import (
	"testing"
	"path"
)

func TestMultipleFileCreateAndListWithMultiFileQueryWithAddSpec(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent := "Test Content"

	attr1 := "foo:hello"
	attr2 := "bar:world"
	attr3 := "flat:fs"

	listDelim := "+"

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
