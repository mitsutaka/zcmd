package sync

const (
	syncPidFile = "/tmp/sync.pid"
)

var optsRsync = []string{"-avP", "--stats", "--delete", "--delete-excluded"}
