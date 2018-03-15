// Package jira implements jira api client
package jira

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Jira is a jira api client
type Jira struct {
	apiURL   string
	auth     Auth
	client   *http.Client
	dumpData bool
}

// Auth holds authentication data for Jira API calls
type Auth struct {
	User string
	Pass string
}

// New returns a Jira API client
func New(apiURL string, auth Auth, dump bool) *Jira {
	return &Jira{
		apiURL:   apiURL,
		auth:     auth,
		client:   &http.Client{},
		dumpData: dump,
	}
}

func (j *Jira) apiGET(path string) ([]byte, error) {

	url := j.apiURL + path // TODO: fix slashes safely
	j.dump("---\nGET Request URL: %q\n---", url)

	bs, err := j.apiRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}

	j.dump("---\n%v\n---", string(bs))
	return bs, nil
}

func (j *Jira) apiPOST(path string, json string) ([]byte, error) {

	url := j.apiURL + path // TODO: fix slashes safely

	j.dump("---\nPOST Request URL: %q\n---", url)

	r := strings.NewReader(json)

	bs, err := j.apiRequest("POST", url, r)
	if err != nil {
		return []byte{}, err
	}

	return bs, nil
}

func (j *Jira) apiRequest(method string, url string, body io.Reader) ([]byte, error) {

	j.dump("URL: %v", url)
	//log("Data: %v", body)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return []byte{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	req.SetBasicAuth(j.auth.User, j.auth.Pass)

	resp, err := j.client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("request failed with status %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

func (j *Jira) dump(f string, args ...interface{}) {
	if j.dumpData {
		log.Printf(f, args...)
	}
}
