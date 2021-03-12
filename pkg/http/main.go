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

// Auth for BasicAuth
type Auth struct {
	Username string
	Password string
}

// Get a BasicAuth authenticated resource
func Get(a Auth, url string) (resp []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return
	}

	client := http.Client{
		Timeout: time.Duration(30 * time.Second),
	}

	req.SetBasicAuth(a.Username, a.Password)
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
func Post(a Auth, url string, data interface{}) (resp []byte, err error) {
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

	req.SetBasicAuth(a.Username, a.Password)
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
