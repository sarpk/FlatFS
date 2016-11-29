// In memory attribute mapper

package FlatFS

import (
	"testing"
)

func TestFileCreateAndRead(t *testing.T) {
	mountPoint := CreateFlatFS()
	Terminate(mountPoint)
}


