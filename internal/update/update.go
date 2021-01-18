package update

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"barf/internal/config"
	"barf/internal/utils"
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
	release, err := getLatestRelease()

	if err != nil {
		return err
	}

	asset := getReleaseAsset(release)

	if asset != nil {
		fmt.Println("Found new version:", release.TagName)

		file, err := ioutil.TempFile("", release.TagName+"-*")

		if err != nil {
			return err
		}

		err = downloadFile(asset.URL, file.Name())

		if err != nil {
			return err
		}

		return upgradeExecutable(file.Name())
	}

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

	fmt.Println(fmt.Sprintf("To upgrade run: sudo mv %s %s", new, name))

	return nil
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
