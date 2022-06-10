package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/shurcooL/githubv4"
)

func GetIssues(data [][]string, variables map[string]interface{}, client *githubv4.Client) ([][]string, error) {
	// If there's existing data, get the end_cursor from the last line to use as the cursor
	if len(data) > 0 {
		variables["cursor"] = githubv4.String(data[len(data)-1][4])
	}

	query := &IssuesQuery{}
	var results [][]string
	var limit ratelimit

	page := 0
	chunk := 100
	for {

		// Check that we can run the query
		if limit.Cost > limit.Remaining {
			return nil, fmt.Errorf("Rate limit error: The query would exceed the rate limit")
		}

		// Run the query
		err := client.Query(context.Background(), query, variables)
		if err != nil {
			return nil, fmt.Errorf("Query error: %v", err)
		}

		if len(query.Repository.Issues.Nodes) == 0 {
			log.Println("No new issues found since last run")
			break
		}

		// Add the data from the query to the output
		for counter, issue := range query.Repository.Issues.Nodes {
			log.Printf("Processing found issue %v", counter+page*chunk)
			results = append(results, []string{
				"false",
				strconv.Itoa(issue.Number),
				issue.Title,
				issue.Url,
				query.Repository.Issues.PageInfo.EndCursor,
			})
		}

		// If this is the last page, break
		if !query.Repository.Issues.PageInfo.HasNextPage {
			break
		}

		// Update the rate limit information
		limit = query.RateLimit

		// Update the cursor to the end cursor of the previous query
		variables["cursor"] = githubv4.String(query.Repository.Issues.PageInfo.EndCursor)

		page++
	}

	return results, nil
}
