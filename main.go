package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"

	"github.com/google/go-github/v69/github"
)

var matchRepoRegex = regexp.MustCompile("github.com/(.*)/(.*)")

func main() {
	var token string
	var repoURL string
	flag.StringVar(&repoURL, "url", "", "Github Repo URL")
	flag.StringVar(&token, "token", "", "Github Token")
	flag.Parse()

	results := matchRepoRegex.FindStringSubmatch(repoURL)
	if len(results) != 3 {
		log.Fatal("not a github repo")
	}

	fmt.Println(New(results[2], results[1], github.NewClient(nil).WithAuthToken(token)).Score())
}
