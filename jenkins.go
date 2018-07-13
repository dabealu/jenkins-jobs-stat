package main

import (
	"encoding/json"
	"fmt"
	"log"
)

var (
	jobsURL     string = "%s/api/json"           // c.Jenkins.Url
	jobURL      string = "%s/job/%s/api/json"    // c.Jenkins.Url, job
	jobBuildURL string = "%s/job/%s/%d/api/json" // c.Jenkins.Url, job, number
)

// jenkins jobs list
type Jobs struct {
	Jobs []Job `json:"jobs"`
}

type Job struct {
	Name string `json:"name"`
}

// get list of jenkins jobs
func listJobs(c *config) []string {
	url := fmt.Sprintf(jobsURL, c.Jenkins.Url)
	code, b := jenkinsGet(url, c.Jenkins.User, c.Jenkins.Password, c.Jenkins.Verify)

	if code != 200 {
		log.Fatalf("List jobs: response code: %d", code)
	}

	jobs := Jobs{}
	res := []string{}

	e := json.Unmarshal(b, &jobs)
	logFatal("List jobs: json", e)

	for _, j := range jobs.Jobs {
		res = append(res, j.Name)
	}

	return res
}

// part of Action struct which contains build variables
type Parameter struct {
	Class string      `json:"_class"`
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

// part of Action struct which contains user login(s)
type Cause struct {
	Class            string `json:"_class"`
	ShortDescription string `json:"shortDescription"`
	UserID           string `json:"userId"`
	UserName         string `json:"userName"`
}

// contains user and build variables
type Action struct {
	Class      string      `json:"_class"`
	RemoteUrls []string    `json:"remoteUrls"`
	Parameters []Parameter `json:"parameters"`
	Causes     []Cause     `json:"causes"`
}

// contains build status, such as FAIL/SUCCESS, timestamp, etc
type GeneralParameters struct {
	Building          bool   `json:"building"`
	Duration          int64  `json:"duration"`
	EstimatedDuration int64  `json:"estimatedDuration"`
	Number            int64  `json:"number"`
	Result            string `json:"result"`
	Timestamp         int64  `json:"timestamp"`
	URL               string `json:"url"`
	BuiltOn           string `json:"builtOn"`
}

// represent all information about single jenkins build
type Build struct {
	Actions []Action `json:"actions"`
	GeneralParameters
}

// get response code and build from jenkins
func getBuild(c *config, job string, number int) (int, Build) {
	url := fmt.Sprintf(jobBuildURL, c.Jenkins.Url, job, number)
	build := Build{}

	code, b := jenkinsGet(url, c.Jenkins.User, c.Jenkins.Password, c.Jenkins.Verify)
	if code != 200 {
		return code, build
	}

	e := json.Unmarshal(b, &build)
	logFatal("Get build: json", e)

	return code, build
}

// list of job build numbers
type BuildNums struct {
	Builds []BuildNum `json:"builds"`
}

type BuildNum struct {
	Number int `json:"number"`
}

// get list of job build numbers
func listBuildNumbers(c *config, job string) []int {
	builds := BuildNums{}
	result := []int{}

	url := fmt.Sprintf(jobURL, c.Jenkins.Url, job)
	code, b := jenkinsGet(url, c.Jenkins.User, c.Jenkins.Password, c.Jenkins.Verify)

	if code != 200 {
		log.Printf("List builds: response code: %d", code)
		return result
	}

	err := json.Unmarshal(b, &builds)
	logFatal("List builds: json", err)

	for _, build := range builds.Builds {
		result = append(result, build.Number)
	}

	return result
}

// server which grabs each job's builds from Jenkins API
// and put them into channel
func getBuilds(c *config, ch chan Build) {
	jobs := []string{}
	numbers := []int{}
	b := Build{}

	for {
		jobs = listJobs(c)

		for _, job := range jobs {
			numbers = listBuildNumbers(c, job) // []int

			for _, num := range numbers {
				_, b = getBuild(c, job, num) // (int, Build)
				ch <- b
			}
		}
		sleep(c.Common.Interval)
	}
}
