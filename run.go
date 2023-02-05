package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Package struct {
	FullName      string
	Description   string
	StarsCount    int
	ForksCount    int
	LastUpdatedBy string
}

func installDeps() {
	// fmt.Println("install deps...")
	file, err := os.Open("requirements.txt")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var count int = 0
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		if len(scanner.Text()) != 0 {
			command := exec.Command("go", "get", "-u", scanner.Text())
			err := command.Run()
			if err == nil {
				count++
			} else {
				fmt.Println("Error installing package: ", scanner.Text())
			}
		}
	}
	fmt.Println(count, "packages installed...")

}

func compile() {
	fmt.Println("compile...")
}

func test() {
	fmt.Println("test...")
}

func main() {
	args := os.Args[1:]
	// fmt.Println(args)
	if args[0] == "install" {
		// installDeps()
	}
	if args[0] == "build" {
		// compile()
	}
	if args[0] == "test" {
		// test()
	}
}

func getnumPR() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "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"},
	)
	tc := oauth2.NewClient(ctx, ts)
	// get go-github client
	client := github.NewClient(tc)

	repo, _, err := client.Repositories.Get(ctx, "Golang-Coach", "Lessons")
	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls", repo)

	resp, errUrl := http.Get(url)
	fmt.Print(resp)
	if errUrl != nil {
		fmt.Println("Error fetching pull requests:", errUrl)
		return
	}

	if err != nil {
		fmt.Printf("Problem in getting repository information %v\n", err)
		os.Exit(1)
	}

	defer resp.Body.Close() //ensure function executes after main function is completed.

	var pullRequests []map[string]interface{} //stores a list of pull requests from API response

	err = json.NewDecoder(resp.Body).Decode(&pullRequests)

	if err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	numPR := len(pullRequests)
	fmt.Println("Number of pull requests:", numPR)

	pack := &Package{
		FullName:    *repo.FullName,
		Description: *repo.Description,
		ForksCount:  *repo.ForksCount,
		StarsCount:  *repo.StargazersCount,
	}

	fmt.Printf("%+v\n", pack)
}
