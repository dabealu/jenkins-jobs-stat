package main

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	datClient *http.Client
	jClient   *http.Client
)

func createHttpClient(verify bool) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !verify},
	}

	return &http.Client{Transport: tr}
}

func logFatal(str string, e error) {
	if e != nil {
		log.Fatalf("%s: %s", str, e)
	}
}

func sleep(seconds int) {
	time.Sleep(time.Second * time.Duration(seconds))
}

// make http GET with basic auth and return status code and body
func jenkinsGet(url, user, password string, verify bool) (int, []byte) {
	req, e := http.NewRequest("GET", url, nil)
	logFatal("Get URL: create request", e)

	if user != "" && password != "" {
		req.SetBasicAuth(user, password)
	}

	resp, e := jClient.Do(req)
	defer resp.Body.Close()
	logFatal("Get URL: do request", e)

	b, e := ioutil.ReadAll(resp.Body)
	logFatal("Get URL: read body", e)

	return resp.StatusCode, b
}

// check whether build already exists in datastore or not
func isBuildExist(url, user, password string, verify bool) bool {
	req, e := http.NewRequest("HEAD", url, nil)
	logFatal("Head URL: create request", e)

	if user != "" && password != "" {
		req.SetBasicAuth(user, password)
	}

	resp, e := datClient.Do(req)
	logFatal("Head URL: do request", e)

	if resp.StatusCode == 200 {
		return true
	}
	return false
}

// save single document in datastore
func createDoc(url, user, password, doc string, verify bool) {
	req, e := http.NewRequest("PUT", url, strings.NewReader(doc))
	logFatal("Create doc: create request", e)

	req.Header.Add("Content-Type", "application/json")

	if user != "" && password != "" {
		req.SetBasicAuth(user, password)
	}

	resp, e := datClient.Do(req)
	logFatal("Create doc: do request", e)

	b, e := ioutil.ReadAll(resp.Body)
	logFatal("Create doc: read response body", e)

	resp.Body.Close()

	if resp.StatusCode > 299 && resp.StatusCode < 200 {
		log.Printf("Create doc: response code not OK: %s", string(b))
	}
}
