package FlatFS

import (
	"testing"
	"github.com/sarpk/FlatFS/flatFS/test/in-memory"
)

var MOUNT_POINT_PATH = FlatFS.GetCurrentDir()

func TestSetupUtilsFunctions(t *testing.T) {
	RecurseThroughFolders("/tmp/lpbckmtpt/", MOUNT_POINT_PATH, t)
}
