package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/v32/github"
)

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
}
