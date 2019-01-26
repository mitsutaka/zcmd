package sync

const (
	syncPidFile = "/tmp/sync.pid"
)

//nolint[gochecknoglobals]
var optsRsync = []string{"-avP", "--stats", "--delete", "--delete-excluded"}
