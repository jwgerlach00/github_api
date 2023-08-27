package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/benoitmasson/plotters/piechart"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

const urlStem = "https://api.github.com/users"

type Repo struct {
	// Important data to decode from JSON received by GET
	ID           int    `json:"id"`
	Name         string `json:"name"`
	LanguagesURL string `json:"languages_url"`
}

func readAccessToken(txtFilepath string) string {
	// Loads GitHub personal access token assumed to be written as a string in txtFilepath
	content, err := os.ReadFile(txtFilepath)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func getAuthRequest(url, accessToken string) *http.Request {
	/// GET request with modified Header to add accessToken. Necessary for using GitHub API without rate limiting
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	return req
}

func getReposWholeResponse(username, accessToken string, client *http.Client) map[string]interface{} {
	// GET request to retrieve the entire response object from /repos using GitHub API
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
	// GET request to retrieve response object according to Repo schema from /repos using GitHub API
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
	// Sums bytes for each language. Returns mapping of language to bytes
	aggregateLanguages := make(map[string]float64)
	for _, repo := range repos {
		languages := getRepoLanguages(repo, accessToken, client)
		for lang, count := range languages {
			aggregateLanguages[lang] += count
		}
	}
	return aggregateLanguages
}

func getRepoLanguages(repo Repo, accessToken string, client *http.Client) map[string]float64 {
	// GET request for language-byte mapping for a repo
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

func plotDistribution(data map[string]float64) {
	// Plots language-byte distribution as a PieChart
	var keys []string
	var values []float64
	for key, value := range data {
		if key != "HTML" && key != "Jupyter Notebook" {
			keys = append(keys, key)
			values = append(values, value)
		}
	}

	p := plot.New()
	p.HideAxes()

	pie, err := piechart.NewPieChart(plotter.Values(values))
	if err != nil {
		log.Fatal(err)
	}
	pie.Labels.Nominal = keys

	pie.Color = color.RGBA{200, 100, 100, 255} // red
	p.Add(pie)
	p.Save(800, 800, "piechart.png")
}

func main() {
	username := "jwgerlach00"
	accessToken := readAccessToken("access_token.txt")

	client := &http.Client{}

	// wholeRes := getReposWholeResponse(username, accessToken, client)

	repos := getRepos(username, accessToken, client)
	aggregateLanguages := countTotalLanguageBytes(repos, accessToken, client)

	fmt.Println(aggregateLanguages)
	plotDistribution(aggregateLanguages)
}
