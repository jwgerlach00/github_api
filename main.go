package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const urlStem = "https://api.github.com/users"

type Repo struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	LanguagesURL string `json:"languages_url"`
}

func getRepos(username string) []Repo {
	url := fmt.Sprintf("%s/%s/repos", urlStem, username)

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close() // close after function resolves

	var repos []Repo
	err = json.NewDecoder(response.Body).Decode(&repos)
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

	repos := getRepos(username)

	for i := 0; i < len(repos); i++ {
		languages := getLanguages(repos[i])
		fmt.Println(languages)
	}
}
