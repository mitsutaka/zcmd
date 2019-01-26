package zcmd

import (
	"reflect"
	"runtime"
	"testing"
)

// GetRsyncCmd returns rsync command and arguments for each platform
func TestGetRsyncCmd(t *testing.T) {
	t.Parallel()

	var expected []string
	expected = append(expected, sudoCmd...)
	switch runtime.GOOS {
	case "linux":
		expected = append(expected, cmdRsyncLinux)
	case "darwin":
		expected = append(expected, cmdRsyncDarwin)
	default:
		t.Fatal("unknown platform")
	}

	cmdRsync, err := GetRsyncCmd()
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(cmdRsync, expected) {
		t.Errorf("%v != %v", cmdRsync, expected)
	}
}
