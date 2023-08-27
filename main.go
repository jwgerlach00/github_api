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

func getReposWholeResponse(username, accessToken string, client *http.Client) map[string]interface{} {
	url := fmt.Sprintf("%s/%s", urlStem, username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

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

func getRepos(username string, client http.Client) []Repo {
	url := fmt.Sprintf("%s/%s/repos", urlStem, username)

	response, err := http.Get(url)
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

func getLanguages(repo Repo) map[string]interface{} {
	response, err := http.Get(repo.LanguagesURL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	var languages map[string]interface{}
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

	getReposWholeResponse(username, accessToken, client)

	// repos := getRepos(username)

	// fmt.Println(repos[0])

	// for i := 0; i < len(repos); i++ {
	// 	languages := getLanguages(repos[i])
	// 	fmt.Println(languages)
	// 	// for key := range languages {
	// 	// 	fmt.Println(key)
	// 	// }
	// }

}
