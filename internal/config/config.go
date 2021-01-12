package config

import (
	"log"
	"os"
	"os/user"
	"path"
)

const DaemonVariable = "_RFT_DAEMON_"
const Name = "rft"
const Description = "A tool for doing robust file operations."

var Version = "v0.0.0"
var BuildTime = "unknown"
var BuildChecksum = "unknown"

var ConfigDir = ""
var PidFile = ""
var LogFile = ""
var SocketFile = ""
var JournalDir = ""

func init() {
	usr, err := user.Current()

	if err != nil {
		log.Fatalf("Could not get user, %+v", err)
	}

	ConfigDir = path.Join(usr.HomeDir, ".config", Name)

	os.Mkdir(ConfigDir, 0700)

	PidFile = path.Join(ConfigDir, "rft.pid")
	LogFile = path.Join(ConfigDir, "rft.log")
	SocketFile = path.Join(ConfigDir, "rft.sock")
	JournalDir = path.Join(ConfigDir, "journal")

	os.Mkdir(JournalDir, 0700)
}
