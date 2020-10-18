package wiremock

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const wiremockAdminURN = "__admin/mappings"

type Client struct {
	url string
}

// NewClient returns *Client.
func NewClient(url string) *Client {
	return &Client{url: url}
}

// StubFor sends http request with StubRule to wiremock server.
func (c *Client) StubFor(stubRule *StubRule) error {
	requestBody, err := json.Marshal(stubRule)
	if err != nil {
		return errors.New(fmt.Sprintf("build stub request error: %s", err.Error()))
	}

	res, err := http.Post(fmt.Sprintf("%s/%s", c.url, wiremockAdminURN), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return errors.New(fmt.Sprintf("stub request error: %s", err.Error()))
	}

	if res.StatusCode != http.StatusCreated {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.New(fmt.Sprintf("read response error: %s", err.Error()))
		}

		return errors.New(fmt.Sprintf("bad response status: %d, response: %s", res.StatusCode, string(bodyBytes)))
	}

	return nil
}

// Clear sends http request to wiremock server for delete all mappings.
func (c *Client) Clear() error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", c.url, wiremockAdminURN), nil)
	if err != nil {
		return errors.New(fmt.Sprintf("build cleare request error: %s", err.Error()))
	}

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return errors.New(fmt.Sprintf("Clear request error: %s", err.Error()))
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("bad response status: %d", res.StatusCode))
	}

	return nil
}