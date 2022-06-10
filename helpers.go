package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	yaml "gopkg.in/yaml.v3"
)

func checkEnvironment() (token, repo, org string, err error) {

	home := os.Getenv("HOME")
	configPath := filepath.Join(home, ".config/gh/hosts.yml")

	ghAuthFileData := make(map[string]githubAuthFile)
	ghAuthFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(ghAuthFile, &ghAuthFileData)
	if err != nil {
		return
	}

	token = ghAuthFileData["github.com"].OAuthToken

	repoFull, ok := os.LookupEnv("GITHUB_REPO")
	if !ok {
		err = fmt.Errorf("Required environment variable GITHUB_REPO not set")
		return
	}

	repoParts := strings.Split(repoFull, "/")
	org = repoParts[0]
	repo = repoParts[1]

	return
}

func createCSV() error {
	_, err := os.Stat(outputFile)

	if os.IsNotExist(err) {
		log.Println("Output file does not exist; creating output file")
		record := []string{"reviewed", "number", "title", "url", "end_cursor"}

		f, err := os.Create(outputFile)
		if err != nil {
			return err
		}
		defer f.Close()

		w := csv.NewWriter(f)

		log.Println("Writing new file headers")
		err = w.Write(record)
		if err != nil {
			return err
		}

		w.Flush()
		return nil
	}

	return nil
}

func readExistingData() ([][]string, error) {
	f, err := os.Open(outputFile)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	r := csv.NewReader(f)

	// skip the header row
	if _, err := r.Read(); err != nil {
		return nil, err
	}

	records, err := r.ReadAll()

	if err != nil {
		return nil, err
	}

	return records, nil
}

func getClient(token string) *githubv4.Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	httpClient := oauth2.NewClient(context.Background(), src)

	return githubv4.NewClient(httpClient)
}

func writeData(data [][]string) error {
	f, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_APPEND, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	err = w.WriteAll(data)
	if err != nil {
		return err
	}

	return nil
}
