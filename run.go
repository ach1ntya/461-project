package main

import (
	"fmt"
	"os"
	"bufio"
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
		installDeps()
	}
	if args[0] == "build" {
		compile()
	}
	if args[0] == "test" {
		test()
	}

}