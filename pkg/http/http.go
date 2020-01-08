package http

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"
)

// HTTPClient represents a http client
func HTTPClient(method, url string) ([]byte, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	client := http.Client{}

	request, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, 0, err
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return body, resp.StatusCode, nil
}
