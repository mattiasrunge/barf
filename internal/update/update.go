package update

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"barf/internal/config"
	"barf/internal/utils"

	"golang.org/x/sys/unix"
)

type githubAsset struct {
	Name string `json:"name"`
	URL  string `json:"browser_download_url"`
}

type githubRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []githubAsset `json:"assets"`
}

// Update updates the binary
func Update() error {
	err := checkWriteAccess()

	if err != nil {
		return err
	}

	release, err := getLatestRelease()

	if err != nil {
		return err
	}

	version := strings.TrimSuffix(release.TagName, "-release")

	if version == config.Version {
		fmt.Println("No new version found")
		return nil

	}

	fmt.Println("Found new version:", version)

	asset := getReleaseAsset(release)

	if asset == nil {
		return errors.New("no files found on release")
	}

	file, err := ioutil.TempFile("", version+"-*")

	if err != nil {
		return err
	}

	err = downloadFile(asset.URL, file.Name())

	if err != nil {
		return err
	}

	err = upgradeExecutable(file.Name())

	if err != nil {
		return err
	}

	name, err := os.Executable()

	if err != nil {
		return err
	}

	fmt.Printf("New version installed at %s\n", name)

	return nil
}

func upgradeExecutable(new string) error {
	err := os.Chmod(new, 0755)

	if err != nil {
		return err
	}

	name, err := os.Executable()

	if err != nil {
		return err
	}

	return os.Rename(new, name)
}

func getReleaseAsset(release *githubRelease) *githubAsset {
	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, config.BuildName) {
			return &asset
		}
	}

	return nil
}

func getLatestRelease() (*githubRelease, error) {
	url := "https://api.github.com/repos/mattiasrunge/barf/releases/latest"
	headers := map[string]string{
		"Accept": "application/vnd.github.v3+json",
	}

	var release githubRelease

	err := utils.GetJSON(url, headers, &release)

	if err != nil {
		return nil, err
	}

	return &release, nil
}

func downloadFile(url string, target string) error {
	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return extractBinary(resp.Body, target)
}

func extractBinary(gzipStream io.Reader, target string) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)

	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if header.Typeflag == tar.TypeReg {
			file, err := os.Create(target)

			if err != nil {
				return err
			}

			defer file.Close()

			_, err = io.Copy(file, tarReader)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func checkWriteAccess() error {
	name, err := os.Executable()

	if err != nil {
		return err
	}

	return unix.Access(name, unix.W_OK)
}
