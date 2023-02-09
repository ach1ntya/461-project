package main

import (
	"context"
	"fmt"
	"regexp"

	"github.com/google/go-github/v32/github"
)

func extractProjectName(inputCommand string) string {
	npmPattern := `npm\s+(install|i)\s+([\w-]+)`
	urlPattern := `(https?://)?(www\.)?([\w-]+)\.([\w-]+)/([\w-]+)/([\w-]+)/?`

	npmRegexp := regexp.MustCompile(npmPattern)
	npmMatch := npmRegexp.FindStringSubmatch(inputCommand)
	if len(npmMatch) > 0 {
		return npmMatch[2]
	}

	urlRegexp := regexp.MustCompile(urlPattern)
	urlMatch := urlRegexp.FindStringSubmatch(inputCommand)
	if len(urlMatch) > 0 {
		return urlMatch[6]
	}

	return ""
}

func main() {
	owner := "joelchiang2k"
	repo := "Linkr"

	client := github.NewClient(nil)

	ctx := context.Background()

	pullRequests, _, err := client.PullRequests.List(ctx, owner, repo, &github.PullRequestListOptions{})
	if err != nil {
		fmt.Println(err)
	}

	pullRequestCount := len(pullRequests)

	fmt.Println("Number of pull requests:", pullRequestCount)

	inputCommand := "npm install express"
	println(extractProjectName(inputCommand))

	inputCommand = "https://github.com/expressjs/express"
	println(extractProjectName(inputCommand))
}
