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
		err = fmt.Errorf("failed to GET (%s) %s", url, err)
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
	if s >= 400 && s < 600 {
		err = fmt.Errorf("received %s", strconv.Itoa(s))
	}
	return
}

// Post an BasicAuth authenticated resource
func (h Client) Post(url string, data interface{}) (resp []byte, err error) {
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

	httpError := hasError(res.StatusCode)

	resp, err = ioutil.ReadAll(res.Body)

	if err != nil {
		return
	}

	if httpError != nil {
		err = fmt.Errorf("failed to POST (%s) %s - %s", url, httpError, string(resp))
		return
	}

	return
}

// PostUnmarshalled makes a POST HTTP request and unmarshalls the data
func (h Client) PostUnmarshalled(url string, data interface{}, target interface{}) (err error) {
	resp, err := h.Post(url, data)

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, target)

	return
}

// GetUnmarshalled makes a GET HTTP request and unmarshalls the data
func (h Client) GetUnmarshalled(url string, targetPtr interface{}) (err error) {
	resp, err := h.Get(url)

	if err != nil {
		err = fmt.Errorf("%s - %s", err.Error(), string(resp))
		return
	}

	err = json.Unmarshal(resp, targetPtr)

	return
}
