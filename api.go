package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type API struct {
	RepoToken       string          `json:"repo_token"`
	RunAt           int64           `json:"run_at"`
	SourceFiles     []*FileCoverage `json:"source_files"`
	CoveredPercent  float64         `json:"covered_percent"`
	CoveredStrength float64         `json:"covered_strength"`
	LineCounts      *LineCounts     `json:"line_counts"`
	Partial         bool            `json:"partial"`
	Git             *Git            `json:"git"`
	Environment     *Environment    `json:"environment"`
	CIService       interface{}     `json:"ci_service"`
}

func (api *API) Post() bool {
	api.validate()

	host := env["CODECLIMATE_API_HOST"]
	if host == "" {
		host = "https://codeclimate.com"
	}

	data, err := json.Marshal(api)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/test_reports", host), bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	if res.StatusCode == 200 {
		fmt.Fprintf(os.Stdout, "%s: %s", res.Status, string(b))

		return true
	}

	fmt.Fprintf(os.Stderr, "%s: %s", res.Status, string(b))

	return false
}

func (api *API) validate() {
	if api.RepoToken == "" {
		panic("repo_token is missing")
	}
	if api.Git == nil {
		panic("git is missing")
	}
	if api.Git.Head == "" {
		panic("git.head is missing")
	}
	if api.Git.CommittedAt == "" {
		panic("git.committed_at is missing")
	}
	if api.SourceFiles == nil {
		panic("source_files are missing")
	}
}
