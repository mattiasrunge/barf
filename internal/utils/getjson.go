package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// GetJSON gets a json from a url
func GetJSON(URL string, headers map[string]string, v interface{}) error {
	spaceClient := http.Client{}

	req, err := http.NewRequest(http.MethodGet, URL, nil)

	if err != nil {
		return err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res, err := spaceClient.Do(req)

	if err != nil {
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}
