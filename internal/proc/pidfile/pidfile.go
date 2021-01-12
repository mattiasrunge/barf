package pidfile

import (
	"io/ioutil"
	"os"
	"strconv"

	"rft/internal/config"
)

func Read() (int, error) {
	if _, err := os.Stat(config.PidFile); err == nil {
		data, _ := ioutil.ReadFile(config.PidFile)

		return strconv.Atoi(string(data))
	} else {
		return 0, err
	}
}

func Write(pid int) error {
	data := []byte(strconv.Itoa(pid))

	return ioutil.WriteFile(config.PidFile, data, 0666)
}

func Delete() error {
	return os.Remove(config.PidFile)
}
