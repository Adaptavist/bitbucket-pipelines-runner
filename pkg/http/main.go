package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Json interface {
	// ToJson encodes the current struct
	ToJson() ([]byte, error)
}

// Auth for BasicAuth
type Auth struct {
	Username string
	Password string
}

// Client for HTTP
type Client struct {
	Auth Auth
}

// NewClient
func NewClient(username string, password string) Client {
	return Client{
		Auth: Auth{
			Username: username,
			Password: password,
		},
	}
}

// Get a BasicAuth authenticated resource
func (h Client) Get(url string) (resp []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return
	}

	client := http.Client{
		Timeout: time.Duration(30 * time.Second),
	}

	req.SetBasicAuth(h.Auth.Username, h.Auth.Password)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		return
	}

	err = hasError(res.StatusCode)

	if err != nil {
		return
	}

	resp, err = ioutil.ReadAll(res.Body)

	if err != nil {
		return
	}

	return
}

// hasError returns an error if 40x or 50x codes are given
func hasError(s int) (err error) {
	// NOTE: Some 400 errors means pipelines are not enabled or a pipelines file doesn't exist!
	if s == 400 || s == 401 || s == 402 || s == 403 || s == 500 || s == 501 || s == 502 || s == 503 {
		err = fmt.Errorf("Received %s", strconv.Itoa(s))
	}
	return
}

// Post an BasicAuth authenticated resource
func (h Client)  Post(url string, data interface{}) (resp []byte, err error) {
	reqData, err := json.Marshal(data)

	if err != nil {
		return
	}

	reqBody := bytes.NewBuffer(reqData)
	req, err := http.NewRequest("POST", url, reqBody)

	if err != nil {
		return
	}

	client := http.Client{
		Timeout: time.Duration(30 * time.Second),
	}

	req.SetBasicAuth(h.Auth.Username, h.Auth.Password)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		return
	}

	err = hasError(res.StatusCode)

	if err != nil {
		return
	}

	resp, err = ioutil.ReadAll(res.Body)

	if err != nil {
		return
	}

	return
}

// PostUnmarshalled makes a POST HTTP request and unmarshalls the data
func (h Client) PostUnmarshalled(url string, data interface{}, targetPtr interface{}) (err error) {
	resp, err := h.Post(url, data)

	if err != nil {
		return fmt.Errorf("post: %s", err)
	}

	err = json.Unmarshal(resp, targetPtr)

	if err != nil {
		return fmt.Errorf("%s: %s", err, string(resp))
	}

	return
}

// GetUnmarshalled makes a GET HTTP request and unmarshalls the data
func (h Client) GetUnmarshalled(url string, targetPtr interface{}) (err error) {
	resp, err := h.Get(url)

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, targetPtr)

	return
}
