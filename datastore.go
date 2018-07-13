package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

// represent single build info which stored in database
type datBuild struct {
	BuildParameters   map[string]interface{} `json:"buildParams"`
	Users             string                 `json:"users"`
	Name              string                 `json:"name"`
	Building          bool                   `json:"building"`
	Duration          int64                  `json:"duration"`
	EstimatedDuration int64                  `json:"estimatedDuration"`
	Number            int64                  `json:"number"`
	Result            string                 `json:"result"`
	Timestamp         string                 `json:"timestamp"`
	URL               string                 `json:"url"`
	BuiltOn           string                 `json:"builtOn"`
}

// create build UID which is used for document's _id in elasticsearch
func (b *datBuild) getID() string {
	data := []byte(b.Timestamp + b.URL)
	return fmt.Sprintf("%x", md5.Sum(data))
}

// extract job name from url
func getJobName(url string) string {
	s := strings.Split(url, "/")
	return s[len(s)-3]
}

// extract build variables from list of actions
func getParams(act []Action) map[string]interface{} {
	bp := make(map[string]interface{})

	for _, a := range act {
		for _, p := range a.Parameters {
			bp[p.Name] = p.Value
		}
	}

	return bp
}

// convert jenkins millisec timestamp into
// string formatted for elasticsearch
func msToStr(ms int64) string {
	t := time.Unix(0, ms*int64(time.Millisecond))
	return t.Format("2006-01-02T15:04:05")
}

// extract list of user(s) from list of actions
func getUsers(act []Action) string {
	users := []string{}

	for _, a := range act {
		for _, c := range a.Causes {
			if len(c.UserID) > 0 {
				users = append(users, c.UserID)
			}
		}
	}

	return strings.Join(users, ";")
}

// convert Build to datBuild which is ready to be stored in DB
func getdatBuild(b Build) datBuild {
	d := datBuild{}

	d.Building = b.Building
	d.Duration = b.Duration
	d.EstimatedDuration = b.EstimatedDuration
	d.Number = b.Number
	d.Result = b.Result
	d.URL = b.URL
	d.BuiltOn = b.BuiltOn

	d.Timestamp = msToStr(b.Timestamp)
	d.Users = getUsers(b.Actions)
	d.BuildParameters = getParams(b.Actions)
	d.Name = getJobName(b.URL)

	return d
}

// server which continuously takes builds from channel
// and saves them to datastore
func saveBuilds(c *config, ch chan Build) {
	b := datBuild{}
	docURL := ""
	db := c.Datastore

	for {
		b = getdatBuild(<-ch)

		docURL = fmt.Sprintf("%s/%s/build/%s", db.Url, db.Index, b.getID())

		// skip build if it doesn't finished yet or
		// if it already exist in datastore
		if b.Building || isBuildExist(docURL, db.User, db.Password, db.Verify) {
			continue
		}

		docBytes, err := json.Marshal(b)
		if err != nil {
			log.Printf("Save builds: json: %v", err)
		}

		createDoc(docURL+"/_create", db.User, db.Password, string(docBytes), db.Verify)
	}
}
