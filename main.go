package main

/*
Script to grab all the repos of the kubernetes github org and compile a list of
the top contributors between all the repos. ENV Variables GITHUB_USER and
GITHUB_TOKEN are used to set Basic Auth.
*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type repoDetails struct {
	Name string `json:"name"`
}

type userDetails struct {
	Login         string `json:"login"`
	Contributions int    `json:"contributions"`
}

func main() {
	org := "kubernetes"
	// Get Repos from orgs
	url := "https://api.github.com/orgs/" + org + "/repos"
	reposList, err := grabRepos(url)
	if err != nil {
		log.Fatal(err)
	}
	// Get Top Contributors
	contributors, err := grabContributors(reposList, org)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(contributors["jbeda"])
	for _, contributor := range contributors {
		fmt.Printf("User: %v\n", contributor.Login)
		fmt.Printf("Contributions: %v\n\n", contributor.Contributions)
	}
}

func grabRepos(url string) ([]repoDetails, error) {
	body, err := processAPI(url)
	if err != nil {
		return nil, err
	}
	reposList := []repoDetails{}
	json.Unmarshal(body, &reposList)
	return reposList, err
}

func grabContributors(reposList []repoDetails, org string) (map[string]userDetails, error) {
	usersMap := make(map[string]userDetails)
	for _, repo := range reposList {
		url := "https://api.github.com/repos/" + org + "/" + repo.Name + "/contributors"
		body, err := processAPI(url)
		if err != nil {
			return nil, err
		}
		users := []userDetails{}
		json.Unmarshal(body, &users)
		usersMap = processUsers(users, usersMap)
	}
	return usersMap, nil
}

func processUsers(users []userDetails, usersMap map[string]userDetails) map[string]userDetails {
	for _, user := range users {
		if val, exist := usersMap[user.Login]; exist {
			val.Contributions += user.Contributions
			usersMap[user.Login] = val
		} else {
			usersMap[user.Login] = userDetails{
				Login:         user.Login,
				Contributions: user.Contributions,
			}
		}
	}
	return usersMap
}

func processAPI(url string) ([]byte, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(os.Getenv("GITHUB_USER"), os.Getenv("GITHUB_TOKEN"))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, err
}
