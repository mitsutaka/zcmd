package zcmd

import (
	"reflect"
	"runtime"
	"testing"
)

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

	cmdRsync, err := getRsyncCmd()
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(cmdRsync, expected) {
		t.Errorf("%v != %v", cmdRsync, expected)
	}
}
