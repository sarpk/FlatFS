package FlatFS

import (
	"testing"
)

func TestUtilsFunctions(t *testing.T) {
	RecurseThroughFolders("/home/sarp/Downloads/impressions-v1/impress_home", "/tmp/mountpoint/", t)
}
