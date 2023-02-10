package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func newURL(url string) *attribute{
	scoreObject := attribute{url: url}
	return &scoreObject
}

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
	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
		scoreObject := newURL(string(line))

		// call sub function to score each line/project
		if strings.Contains(line, "github.com") {
			github(line, scoreObject)
		} else if strings.Contains(line, "npmjs.com") {
			npmjs(line, scoreObject)
		} else {
			fmt.Println("Error: ", line, "is not a valid URL")
		}
	}
}

/*func scoreGitHub(scoreObject attribute){
	print("no score yet")
}

func scoreNPM(scoreObject *attribute) {
	npmRestAPI(scoreObject)
	fmt.Println("url:", scoreObject.url)
	fmt.Println("\nMaintainers:", scoreObject.responsiveness)
	fmt.Println("\nlicense:", scoreObject.license)
	return
}*/

func github(url string, scoreOject *attribute) {
	split := strings.Split(url, "/")
	owner := split[len(split)-2]
	repo := split[len(split)-1]
	print("Owner: ", owner, " Repo: ", repo, "\n")
}

func npmjs(url string, scoreObject *attribute) {
	split := strings.Split(url, "/")
	packageName := split[len(split)-1]
	print("Package: ", packageName, "\n")
	npmRestAPI(packageName, scoreObject)
}

//func npmRestAPI(packageName string) {

//	url := "https://registry.npmjs.org/" + packageName

func npmRestAPI(packageName string, scoreObject *attribute) {
	
	//http get request to connect to registry api of package
	//response, err := http.Get("https://registry.npmjs.org/express")
	url := "https://registry.npmjs.org/" + packageName
	response, err := http.Get(url)

	//update rest API method to dynamically take in package name

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
	fmt.Print("\n")

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
