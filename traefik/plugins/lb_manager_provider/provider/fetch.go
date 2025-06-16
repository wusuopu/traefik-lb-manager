package provider

import (
	"fmt"
	"io"
	"net/http"
	"time"
)


func request(url string, method string) ([]byte, error) {
  client := &http.Client{
		Timeout: 3 * time.Second,
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode > 300 {
		return nil, fmt.Errorf("request %s %s failed with status code %d", method, url, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
