package config

import (
	"log"
	"os"
	"os/user"
	"path"
)

const DaemonVariable = "_barf_DAEMON_"
const Name = "barf"
const Description = "A tool for doing robust file operations."

var Version = "v0.0.0"
var BuildTime = "unknown"
var BuildChecksum = "unknown"
var production = "no"

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

	parentDir := path.Join(usr.HomeDir, ".config")

	os.Mkdir(parentDir, 0700)

	ConfigDir = path.Join(parentDir, Name)

	os.Mkdir(ConfigDir, 0700)

	PidFile = path.Join(ConfigDir, "barf.pid")
	LogFile = path.Join(ConfigDir, "barf.log")
	SocketFile = path.Join(ConfigDir, "barf.sock")
	JournalDir = path.Join(ConfigDir, "journal")

	os.Mkdir(JournalDir, 0700)
}

// IsProduction returns true if production build
func IsProduction() bool {
	return production != "no"
}
