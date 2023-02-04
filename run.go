package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

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

func main() {
	args := os.Args[1:]
	if args[0] == "install" {
		installDeps()
	} else if args[0] == "build" {
		compile()
	} else if args[0] == "test" {
		test()
	} else if args[0] == "URL FILE" {
		//add function and function call here
	} else {
		help()
		os.Exit(1)
	}

}
