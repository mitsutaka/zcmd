package nas

const (
	nasPidFile = "/tmp/nas.pid"
)

var optsRsync = []string{"-avP", "--stats", "--delete", "--delete-excluded"}
