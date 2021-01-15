package fs

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

func ensureDirectoryFor(filename string) error {
	dirname := path.Dir(filename)
	return os.MkdirAll(dirname, 0700)
}

func writeTempFile(filename string, bytes []byte) error {
	tmpfilename := filename + "~"

	return ioutil.WriteFile(tmpfilename, bytes, 0600)
}

func moveTempFile(filename string) error {
	tmpfilename := filename + "~"

	stat, err := os.Stat(tmpfilename)

	if err == nil {
		if stat.Size() > 0 {
			return os.Rename(tmpfilename, filename)
		}

		return os.Remove(tmpfilename)
	}

	return nil
}

// WriteJSONFile writes the supplied structure as JSON via a temporary file for safety
func WriteJSONFile(filename string, b interface{}) error {
	err := ensureDirectoryFor(filename)

	if err != nil {
		return err
	}

	bytes, err := json.Marshal(b)

	if err != nil {
		return err
	}

	err = writeTempFile(filename, bytes)

	if err != nil {
		return err
	}

	return moveTempFile(filename)
}

// ReadJSONFile reads a JSON file into the supplied structure, first checks if there is a temporary file
func ReadJSONFile(filename string, b interface{}) error {
	err := ensureDirectoryFor(filename)

	if err != nil {
		return err
	}

	err = moveTempFile(filename)

	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return err
	}

	return json.Unmarshal(data, b)
}
