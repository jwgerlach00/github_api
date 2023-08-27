package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const urlStem = "https://api.github.com/users"

type Repo struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	LanguagesURL string `json:"languages_url"`
}

func readAccessToken(filepath string) string {
	content, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func getAuthRequest(url, accessToken string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	return req
}

func getReposWholeResponse(username, accessToken string, client *http.Client) map[string]interface{} {
	url := fmt.Sprintf("%s/%s", urlStem, username)

	req := getAuthRequest(url, accessToken)

	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	responseBytes, err := io.ReadAll(response.Body)
	var wholeResponse map[string]interface{}
	err = json.Unmarshal(responseBytes, &wholeResponse)
	if err != nil {
		log.Fatal(err)
	}

	return wholeResponse
}

func getRepos(username, accessToken string, client *http.Client) []Repo {
	url := fmt.Sprintf("%s/%s/repos", urlStem, username)

	req := getAuthRequest(url, accessToken)

	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close() // close after function resolves

	var repos []Repo
	err = json.NewDecoder(response.Body).Decode(&repos) // gives address of repos rather than repos itself. '*' would
	// give value, not address
	if err != nil {
		log.Fatal(err)
	}

	return repos
}

func countTotalLanguageBytes(repos []Repo, accessToken string, client *http.Client) map[string]float64 {
	aggregateLanguages := make(map[string]float64)
	for _, repo := range repos {
		languages := getLanguages(repo, accessToken, client)
		for lang, count := range languages {
			aggregateLanguages[lang] += count
		}
	}
	return aggregateLanguages
}

func getLanguages(repo Repo, accessToken string, client *http.Client) map[string]float64 {
	req := getAuthRequest(repo.LanguagesURL, accessToken)

	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	var languages map[string]float64
	err = json.NewDecoder(response.Body).Decode(&languages)
	if err != nil {
		log.Fatal(err)
	}

	return languages
}

func main() {
	username := "jwgerlach00"
	accessToken := readAccessToken("access_token.txt")

	client := &http.Client{}

	// wholeRes := getReposWholeResponse(username, accessToken, client)

	repos := getRepos(username, accessToken, client)
	aggregateLanguages := countTotalLanguageBytes(repos, accessToken, client)

	fmt.Println(aggregateLanguages)
}
