package main

import (
	"log"
	"time"

	"github.com/shurcooL/githubv4"
)

const outputFile = "issues_report.csv"

func main() {
	// Check for required env vars
	token, repo, org, err := checkEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	// Create the output file if it does not exist
	if err := createCSV(); err != nil {
		log.Fatal(err)
	}

	// Get the existing data
	data, err := readExistingData()
	if err != nil {
		log.Fatal(err)
	}

	// Get a GitHub client
	client := getClient(token)

	// Make variables
	variables := map[string]interface{}{
		"owner":      githubv4.String(org),
		"repository": githubv4.String(repo),
		// One year before query data
		"start":  githubv4.DateTime{time.Now().AddDate(-1, 0, 0)},
		"cursor": (*githubv4.String)(nil),
	}

	results, err := GetIssues(data, variables, client)
	if err != nil {
		log.Fatal(err)
	}

	if err := writeData(results); err != nil {
		log.Fatal(err)
	}
}
