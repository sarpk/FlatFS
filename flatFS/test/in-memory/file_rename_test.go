// In memory attribute mapper

package FlatFS

import (
	"testing"
	"path"
)

func TestSingleFileRename(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent := "Test Content"

	attr1 := "foo:hello"
	attr2 := "bar:world"
	attr3 := "flat:fs"

	initialFilePath := path.Join(mountPoint, attr1 + "," + attr2)
	movedFilePath := path.Join(mountPoint, attr3 + "," + attr2)

	write_to_file(initialFilePath, testContent)
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

	attr1 := "foo:hello"
	attr2 := "bar:world"
	attr3 := "flat:fs"

	exactPath1 := path.Join(mountPoint, attr1 + "," + attr2)
	exactPath2 := path.Join(mountPoint, attr2 + "," + attr3)

	write_to_file(exactPath1, testContent1)
	write_to_file(exactPath2, testContent2)
	exec_cmd("mv " + exactPath1 + " " + exactPath2)
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

func TestRenameExistingFileToASpec(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent := "Test Content"

	attr1 := "foo:hello"
	attr2 := "bar:world"
	attr3 := "flat:fs"

	addSpec := "?"

	exactPath := path.Join(mountPoint, attr1 + "," + attr2)
	newPath := path.Join(mountPoint, attr1 + "," + attr2 + "," + attr3)
	specPath := path.Join(mountPoint, addSpec + attr3)

	write_to_file(exactPath, testContent)
	exec_cmd("mv " + exactPath + " " + specPath)
	fileContent := read_from_file(newPath)
	assert_string_equals(fileContent, testContent, t)

	fileContent = read_from_file(exactPath)
	assert_string_equals(fileContent, "", t) //File doesn't exist

	lsContent := exec_cmd("ls -l " + newPath)
	lsContent = assert_string_contains_per_line(lsContent, []string{attr1, attr2, attr3}, t)
	assert_string_not_contains(lsContent, attr1, t)
	assert_string_not_contains(lsContent, attr2, t)
	assert_string_not_contains(lsContent, attr3, t)

	Terminate(mountPoint)
}

func TestRenameExistingFileToASpecWhereAFileExists(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent1 := "Test Content"
	testContent2 := "This file should no longer exist"

	attr1 := "foo:hello"
	attr2 := "bar:world"
	attr3 := "flat:fs"

	addSpec := "?"

	exactPath := path.Join(mountPoint, attr1 + "," + attr2)
	newPath := path.Join(mountPoint, attr1 + "," + attr2 + "," + attr3)
	newPathDiffOrder := path.Join(mountPoint, attr3 + "," + attr2 + "," + attr1)
	specPath := path.Join(mountPoint, addSpec + attr3)

	write_to_file(exactPath, testContent1)
	write_to_file(newPath, testContent2)
	fileContent := read_from_file(newPath)
	assert_string_equals(fileContent, testContent2, t)
	exec_cmd("mv " + exactPath + " " + specPath)
	fileContent = read_from_file(newPathDiffOrder)
	assert_string_equals(fileContent, testContent1, t)

	fileContent = read_from_file(exactPath)
	assert_string_equals(fileContent, "", t) //File doesn't exist

	lsContent := exec_cmd("ls -l " + newPath)
	lsContent = assert_string_contains_per_line(lsContent, []string{attr1, attr2, attr3}, t)
	assert_string_not_contains(lsContent, attr1, t)
	assert_string_not_contains(lsContent, attr2, t)
	assert_string_not_contains(lsContent, attr3, t)

	Terminate(mountPoint)
}

func TestRenameExistingQueryWithAReplaceSpec(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent1 := "Test Content"

	attr1 := "foo:hello"
	attr2 := "bar:world"
	attr3 := "flat:fs"

	addSpec := "?"
	replaceSpec := "="

	exactPath := path.Join(mountPoint, attr1 + "," + attr2)
	exactPathDiffOrder := path.Join(mountPoint, attr2 + "," + attr1)
	newPath := path.Join(mountPoint, attr1 + "," + attr3)
	addPath := path.Join(mountPoint, addSpec + attr2)
	replacePath := path.Join(mountPoint, replaceSpec + attr3)

	write_to_file(exactPath, testContent1)
	exec_cmd("mv " + addPath + " " + replacePath)
	fileContent := read_from_file(newPath)
	assert_string_equals(fileContent, testContent1, t)

	fileContent = read_from_file(exactPathDiffOrder)
	assert_string_equals(fileContent, "", t) //File doesn't exist

	lsContent := exec_cmd("ls -l " + newPath)
	lsContent = assert_string_contains_per_line(lsContent, []string{attr1, attr3}, t)
	assert_string_not_contains(lsContent, attr1, t)
	assert_string_not_contains(lsContent, attr2, t)
	assert_string_not_contains(lsContent, attr3, t)

	Terminate(mountPoint)
}

func TestRenameExistingQueryWithAReplaceSpecWhenFileExists(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent1 := "Test Content"
	testContent2 := "This file should no longer exist"

	attr1 := "foo:hello"
	attr2 := "bar:world"
	attr3 := "flat:fs"

	addSpec := "?"
	replaceSpec := "="

	exactPath := path.Join(mountPoint, attr1 + "," + attr2)
	exactPathDiffOrder := path.Join(mountPoint, attr2 + "," + attr1)
	newPath := path.Join(mountPoint, attr1 + "," + attr3)
	newPathDiffOrder := path.Join(mountPoint, attr3 + "," + attr1)
	addPath := path.Join(mountPoint, addSpec + attr2)
	replacePath := path.Join(mountPoint, replaceSpec + attr3)

	write_to_file(exactPath, testContent1)
	write_to_file(newPath, testContent2)
	fileContent := read_from_file(newPath)
	assert_string_equals(fileContent, testContent2, t)
	exec_cmd("mv " + addPath + " " + replacePath)
	fileContent = read_from_file(newPathDiffOrder)
	assert_string_equals(fileContent, testContent1, t)

	fileContent = read_from_file(exactPathDiffOrder)
	assert_string_equals(fileContent, "", t) //File doesn't exist

	lsContent := exec_cmd("ls -l " + newPath)
	lsContent = assert_string_contains_per_line(lsContent, []string{attr1, attr3}, t)
	assert_string_not_contains(lsContent, attr1, t)
	assert_string_not_contains(lsContent, attr2, t)
	assert_string_not_contains(lsContent, attr3, t)

	Terminate(mountPoint)
}

func TestRenameExistingQueryWithADeleteSpec(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent1 := "Test Content"

	attr1 := "foo:hello"
	attr2 := "bar:world"
	attr3 := "flat:fs"

	addSpec := "?"
	deleteSpec := "-"

	exactPath := path.Join(mountPoint, attr1 + "," + attr2 + "," + attr3)
	exactPathDiffOrder := path.Join(mountPoint, attr3 + "," + attr2 + "," + attr1)
	newPath := path.Join(mountPoint, attr1 + "," + attr2)
	addPath := path.Join(mountPoint, addSpec + attr2)
	deletePath := path.Join(mountPoint, deleteSpec + attr3)

	write_to_file(exactPath, testContent1)
	exec_cmd("mv " + addPath + " " + deletePath)
	fileContent := read_from_file(newPath)
	assert_string_equals(fileContent, testContent1, t)

	fileContent = read_from_file(exactPathDiffOrder)
	assert_string_equals(fileContent, "", t) //File doesn't exist

	lsContent := exec_cmd("ls -l " + newPath)
	lsContent = assert_string_contains_per_line(lsContent, []string{attr1, attr2}, t)
	assert_string_not_contains(lsContent, attr1, t)
	assert_string_not_contains(lsContent, attr2, t)
	assert_string_not_contains(lsContent, attr3, t)

	Terminate(mountPoint)
}

func TestRenameExistingQueryWithADeleteSpecWhenFileExists(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent1 := "Test Content"
	testContent2 := "This file should no longer exist"

	attr1 := "foo:hello"
	attr2 := "bar:world"
	attr3 := "flat:fs"

	addSpec := "?"
	deleteSpec := "-"

	exactPath := path.Join(mountPoint, attr1 + "," + attr2 + "," + attr3)
	exactPathDiffOrder := path.Join(mountPoint, attr2 + "," + attr1 + "," + attr3)
	newPath := path.Join(mountPoint, attr1 + "," + attr2)
	newPathDiffOrder := path.Join(mountPoint, attr2 + "," + attr1)
	addPath := path.Join(mountPoint, addSpec + attr2)
	deletePath := path.Join(mountPoint, deleteSpec + attr3)

	write_to_file(exactPath, testContent1)
	write_to_file(newPath, testContent2)
	fileContent := read_from_file(newPath)
	assert_string_equals(fileContent, testContent2, t)
	exec_cmd("mv " + addPath + " " + deletePath)
	fileContent = read_from_file(newPathDiffOrder)
	assert_string_equals(fileContent, testContent1, t)

	fileContent = read_from_file(exactPathDiffOrder)
	assert_string_equals(fileContent, "", t) //File doesn't exist

	lsContent := exec_cmd("ls -l " + newPath)
	lsContent = assert_string_contains_per_line(lsContent, []string{attr1, attr2}, t)
	assert_string_not_contains(lsContent, attr1, t)
	assert_string_not_contains(lsContent, attr2, t)
	assert_string_not_contains(lsContent, attr3, t)

	Terminate(mountPoint)
}


func TestRenameExistingQueryWithAnAddSpec(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent1 := "Test Content"

	attr1 := "foo:hello"
	attr2 := "bar:world"
	attr3 := "flat:fs"

	addSpec := "?"
	querySpec := "?"

	exactPath := path.Join(mountPoint, attr1 + "," + attr2)
	exactPathDiffOrder := path.Join(mountPoint, attr2 + "," + attr1)
	newPath := path.Join(mountPoint, attr1 + "," + attr2 + "," + attr3)
	addPath := path.Join(mountPoint, addSpec + attr2)
	queryPath := path.Join(mountPoint, querySpec + attr3)

	write_to_file(exactPath, testContent1)
	exec_cmd("mv " + addPath + " " + queryPath)
	fileContent := read_from_file(newPath)
	assert_string_equals(fileContent, testContent1, t)

	fileContent = read_from_file(exactPathDiffOrder)
	assert_string_equals(fileContent, "", t) //File doesn't exist

	lsContent := exec_cmd("ls -l " + newPath)
	lsContent = assert_string_contains_per_line(lsContent, []string{attr1, attr2, attr3}, t)
	assert_string_not_contains(lsContent, attr1, t)
	assert_string_not_contains(lsContent, attr2, t)
	assert_string_not_contains(lsContent, attr3, t)

	Terminate(mountPoint)
}

func TestRenameExistingQueryWithAnAddSpecWithSpecChanged(t *testing.T) {
	mountPoint := CreateFlatFS()
	testContent1 := "Test Content"

	attr1 := "foo:hello"
	attr2 := "bar:world"
	attr3 := "flat:fs"

	addSpec := "+"
	querySpec := "+"

	exactPath := path.Join(mountPoint, attr1 + "," + attr2)
	exactPathDiffOrder := path.Join(mountPoint, attr2 + "," + attr1)
	newPath := path.Join(mountPoint, attr1 + "," + attr2 + "," + attr3)
	addPath := path.Join(mountPoint, addSpec + attr2)
	queryPath := path.Join(mountPoint, querySpec + attr3)

	write_to_file(exactPath, testContent1)
	exec_cmd("mv " + addPath + " " + queryPath)
	fileContent := read_from_file(newPath)
	assert_string_equals(fileContent, testContent1, t)

	fileContent = read_from_file(exactPathDiffOrder)
	assert_string_equals(fileContent, "", t) //File doesn't exist

	lsContent := exec_cmd("ls -l " + newPath)
	lsContent = assert_string_contains_per_line(lsContent, []string{attr1, attr2, attr3}, t)
	assert_string_not_contains(lsContent, attr1, t)
	assert_string_not_contains(lsContent, attr2, t)
	assert_string_not_contains(lsContent, attr3, t)

	Terminate(mountPoint)
}
