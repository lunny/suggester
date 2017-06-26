package suggester

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Result describes the returned
type Result struct {
	ID   int64
	Item string
}

// Client describes the client connecting web service
type Client struct {
	URLPrefix string
}

// New creates a new client
func New(urlPrefix string) *Client {
	return &Client{strings.TrimRight(urlPrefix, "/")}
}

// AddIndex add some indexes
func (c *Client) AddIndex(prefix, word string, id, unitID int64) error {
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/%s/%d/%s/%d", c.URLPrefix, prefix, unitID, word, id), nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result = make(map[string]string)
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if resp.StatusCode/100 != 2 {
		return errors.New(result["err"])
	}

	return nil
}

// DelIndex delete some indexes
func (c *Client) DelIndex(prefix, word string, id, unitID int64) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/%d/%s/%d", c.URLPrefix, prefix, unitID, word, id), nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result = make(map[string]string)
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if resp.StatusCode/100 != 2 {
		return errors.New(result["err"])
	}

	return nil
}

type response struct {
	Results []Result `json:"results"`
}

// Search search the result
func (c *Client) Search(prefix, kw string, unitID int64, limit int) ([]Result, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%d/%s?limit=%d", c.URLPrefix, prefix, unitID, kw, limit), nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		var result = make(map[string]string)
		if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, err
		}
		return nil, errors.New(result["err"])
	}

	var response response
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Results, nil
}
