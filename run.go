package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"strconv"
	"encoding/binary"
	"math"
	//"github.com/google/go-github/v50/github"
	"github.com/machinebox/graphql"
	//"golang.org/x/oauth2"
)

type attribute struct{
	url string
	netScore string
	rampUp float32
	correctness float32	
	busFactor float32	
	responsiveness float32
	license string 
}

type gitObject struct{
	numCommits float32
	numPullRequests float32
	graphQL float32
}

type npmObject struct{
	numCommits float32
	numMaintainers float32
	graphQL float32
	gitRepo string
}

func newURL(url string) (*attribute) {
	scoreObject := attribute{url: url}
	return &scoreObject
}
/*func newNpmObject(url string) (*npmObject) {
	npmObj := npmObject{}
	return &npmObj
}
func newGitObject(url string) (*gitObject) {
	gitObj := {url: url}
	return &npmObj
}*/

func installDeps() {
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
	if count == 0 {
		fmt.Println("No packages installed.")
		os.Exit(1)
	} else {
		fmt.Println(count, "packages installed...")
		os.Exit(0)
	}
}

func compile() {
	command := exec.Command("go", "build", "run.go")
	err := command.Run()
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	} else {
		fmt.Println("Build successful...")
		os.Exit(0)
	}
}

func test() {
	fmt.Println("test...")
}

func help() {
	fmt.Println("Unknown command\nUsage: ./run [command] [args]\nCommands:\n\tinstall\t\tInstall dependencies\n\tbuild\t\tBuild the project\n\ttest\t\tRun tests\n\tURL FILE\tScore all URLs in the file")
}

func file(filename string) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var	urlCount int = 0
	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
		scoreObject := newURL(string(line))

		// call sub function to score each line/project
		if strings.Contains(line, "github.com") {
			//var gitObj gitObject
			urlCount += 1
			githubFunc(line, scoreObject, urlCount)
		} else if strings.Contains(line, "npmjs.com") {
			var npmObj npmObject
			urlCount += 1
			npmjs(line, scoreObject, urlCount, &npmObj)
		} else {
			fmt.Println("Error: ", line, "is not a valid URL")
		}
	}
}


func githubFunc(url string, scoreObject *attribute, count int) {
	split := strings.Split(url, "/")
	owner := split[len(split)-2]
	repo := split[len(split)-1]
	print("Owner: ", owner, " Repo: ", repo, "\n")
	value1 := githubSource(scoreObject, url, count)
	var fullRepo string = owner + "/" + repo
	value2 := githubPullReq(fullRepo, scoreObject, url)
	//value3 := githubGraphQL
	//intConv, _ := strconv.Atoi(string(value))
  	//gitHubGraphQL(repo, owner)
	fmt.Println(value1)
	fmt.Println(value2)

/*func githubFunc(url string) {
	split := strings.Split(url, "/")
	owner := split[len(split)-2]
	repo := split[len(split)-1]
	// print("Owner: ", owner, " Repo: ", repo, "\n")
	gitHubGraphQL(repo, owner)
	gitHubRestAPI(repo, owner)


}*/
}

func npmjs(url string, scoreObject *attribute, count int, npmObj *npmObject) {
	split := strings.Split(url, "/")
	packageName := split[len(split)-1]
	print("Package: ", packageName, "\n")
	npmRestAPI(packageName, scoreObject, npmObj)
	//npmGraphQL
	npmSource(npmObj, count)
	fmt.Println(npmObj.gitRepo)
	fmt.Println(npmObj.numCommits)
}

func githubSource(scoreObject *attribute, url string, count int) (output []byte){

	//call python script that clones repo and pull number of commits
	command := exec.Command("python3", "cloner.py", url, strconv.Itoa(count))
	output, err := command.Output()

	if err != nil{
		fmt.Println(err.Error())
		return
	}
	
	return output

}

func npmSource(npmObj *npmObject, count int) {

	//call python script that clones repo and pull number of commits
	command := exec.Command("python3", "cloner.py", "https://github.com/" + npmObj.gitRepo, strconv.Itoa(count))
	output, err := command.Output()

	if err != nil{
		fmt.Println(err.Error())
		return
	}
	b := binary.BigEndian.Uint32(output)
	float := math.Float32frombits(b)
	npmObj.numCommits = float

}

/*type PullRequests struct {
	TotalCount int `json:"total_count"`
}*/

func githubPullReq(repoName string, scoreObject *attribute, url string) (value2 int) {
	req, _ := http.NewRequest("GET", "https://api.github.com/search/issues?q=is:pr+repo:" + repoName, nil)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	defer res.Body.Close()

	var pullRequests PullRequests
	json.NewDecoder(res.Body).Decode(&pullRequests)

	return pullRequests.TotalCount
	//fmt.Printf("Number of pull requests in joelchiang2k/Linkr: %d\n", pullRequests.TotalCount)
}

func npmRestAPI(packageName string, scoreObject *attribute, npmObj *npmObject) {
	
	//append packageName to the api url and send request
	url := "https://registry.npmjs.org/" + packageName
	response, err := http.Get(url)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	//read api response into responseData
	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
		//log.Fatal(err)
	}

	//creates an object to store json data
	contributors := make(map[string]interface{})

	//Unmarshalls json data into contributors object and returns err
	err = json.Unmarshal(responseData, &contributors)

	if err != nil {
		fmt.Print("failed to decode api response: ", err)
		os.Exit(1)
	}

	//stores list of maintainers into array object
	array := contributors["maintainers"].([]interface{})
	numContributors := len(array) //number of active maintainers for package
	scoreObject.responsiveness = float32(numContributors)
	license := contributors["license"]
	scoreObject.license = license.(string)
	//fmt.Print(scoreObject.responsiveness)


	//output numContributors and license

	fmt.Print("number of contributors: ", scoreObject.responsiveness)
	fmt.Print("\nlicense: ", contributors["license"].(string))
	//fmt.Print("\ngithub url: ", contributors["repository"].(map[string]interface{})["url"])
	split := strings.Split(contributors["repository"].(map[string]interface{})["url"].(string), "/")
	owner := split[len(split)-2]
	repo := split[len(split)-1]
	fmt.Println("Owner: " ,owner, " Repo: ", repo)
	//var npmGitString string = owner + "/" + repo
	npmObj.gitRepo = owner + "/" + repo
	//fmt.Print("\ngithub url: ", contributors["repository"].(map[string]interface{})["url"])
	fmt.Print("\n")

}

func licenseCompatability(license string) (compatible bool) {
	licenseArr := [6]string{"MIT", "X11", "Public Domain", "BSD-new", "Apache 2.0", "LGPLv2.1"}
	
	for _, l := range licenseArr{
		if l == license{
			return true
		}
	}

	return false

	//fmt.Print("number of contributors: ", numContributors)
	//fmt.Print("\nlicense: ", contributors["license"])
}

/*func gitHubRestAPI(repo string, owner string) {
	apiKey := os.Getenv("GITHUB_API_KEY")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: apiKey},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	_, _, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		fmt.Println(err)
		return
	}

	pullRequests, _, err := client.PullRequests.List(ctx, owner, repo, nil)
	if err != nil {
		fmt.Println("Error fetching pull requests:", err)
		return
	}

	totalPullRequests := len(pullRequests)
	fmt.Printf("Total pull requests: %d\n", totalPullRequests)

	totalIssues, _, err := client.Issues.ListByRepo(ctx, owner, repo, nil)
	fmt.Printf("Total issues: %d\n", len(totalIssues))
}*/

type PullRequests struct {
	TotalCount int `json:"total_count"`
}

func numPullReq() {
	url := "https://api.github.com/search/issues?q=is:pr+repo:joelchiang2k/Linkr" //change repo name OWNER/REPO
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	res, _ := client.Do(req)

	defer res.Body.Close()

	var pullRequests PullRequests
	json.NewDecoder(res.Body).Decode(&pullRequests)

	fmt.Printf("Number of pull requests in joelchiang2k/Linkr: %d\n", pullRequests.TotalCount)
}

func gitHubGraphQL(repoName string, owner string) {
	client := graphql.NewClient("https://api.github.com/graphql")
	req := graphql.NewRequest(`
	query ($repoName: String!, $owner: String!) {
		repository(name: $repoName, owner: $owner) {
		  licenseInfo {
			name
		  }
		  pullRequests {
			totalCount
		  }
		  commitComments {
			totalCount
		  }
		  releases {
			totalCount
		  }
		  stargazerCount
		  defaultBranchRef {
			name
			target {
			  ... on Commit {
				id
				history(first: 0) {
				  totalCount
				}
			  }
			}
		  }
		}
	  }  
	`)
	req.Var("repoName", repoName)
	req.Var("owner", owner)
	apiKey := os.Getenv("GITHUB_API_KEY")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	var response map[string]interface{}
	err := client.Run(context.Background(), req, &response)
	if err != nil {
		panic(err)
	}
	repository := response["repository"].(map[string]interface{})
	commitCount := int(repository["defaultBranchRef"].(map[string]interface{})["target"].(map[string]interface{})["history"].(map[string]interface{})["totalCount"].(float64))
	pullRequests := int(repository["pullRequests"].(map[string]interface{})["totalCount"].(float64))
	releases := int(repository["releases"].(map[string]interface{})["totalCount"].(float64))
	stargazerCount := int(repository["stargazerCount"].(float64))
	fmt.Println("Commit Count: ", commitCount)
	fmt.Println("Pull Requests: ", pullRequests)
	fmt.Println("Releases: ", releases)
	fmt.Println("Stargazers: ", stargazerCount)

}

func main() {
	args := os.Args[1:]
	if args[0] == "install" {
		installDeps()
	} else if args[0] == "build" {
		compile()
	} else if args[0] == "test" {
		test()
	} else if filepath.Ext(args[0]) == ".txt" {
		file(args[0])
	} else {
		help()
		os.Exit(1)
	}
}
