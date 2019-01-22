package nas

const (
	nasPidFile  = "/tmp/nas.pid"
	syncAllPath = "all"
)

var optsRsync = []string{"-avP", "--stats", "--delete", "--delete-excluded"}
