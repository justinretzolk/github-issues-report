package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shurcooL/githubv4"
)

const outputFile = "issues_report.csv"

func main() {
	// Check for required env vars
	token, repo, err := checkEnvironment()
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

	// Get the time for the search query, one year before current time
	now := time.Now().AddDate(-1, 0, 0).Format("2006-01-02")

	// Make variables
	variables := map[string]interface{}{
		"searchQuery": githubv4.String(fmt.Sprintf("repo:%s, is:open, is:issue, created:>%s, sort:created-asc", repo, now)),
		"searchType": githubv4.SearchTypeIssue,
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
