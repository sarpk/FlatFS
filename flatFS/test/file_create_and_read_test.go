// In memory attribute mapper

package FlatFS

import (
	"testing"
	"path"
)

func TestFileCreateAndRead(t *testing.T) {
	mountPoint := CreateFlatFS()
	write_to_file(path.Join(mountPoint, "foo=hello,bar=world"), "test")
	exec_cmd("ls -l " + path.Join(mountPoint, "?foo=hello"))
	read_from_file(path.Join(mountPoint, "foo=hello,bar=world"))
	//TODO assert these two above
	Terminate(mountPoint)
}


