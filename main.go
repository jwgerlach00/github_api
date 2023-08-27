package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const urlStem = "https://api.github.com/users"

type Repo struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	LanguagesURL string `json:"languages_url"`
}

func getReposWholeResponse(username string) {
	url := fmt.Sprintf("%s/%s/repos", urlStem, username)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	responseBytes, err := io.ReadAll(response.Body)
	fmt.Println(string(responseBytes))
}

func getRepos(username string) []Repo {
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
	// accessToken := ""
	getReposWholeResponse(username)

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
