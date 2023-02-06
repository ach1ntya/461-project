package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	folder := "/Users/dkholode/Desktop/ECE 461/461-project/" // Cloned Repo folder

	// Git cmd for list of all repos
	out, err := exec.Command("git", "-C", folder, "branch", "-a").Output()
	if err != nil {
		fmt.Println("Error running git command:", err)
		return
	}

	// Split output onto new lines and return (len - extra versions of origin/head)
	branches := strings.Split(string(out), "\n")
	fmt.Printf("%d", len(branches)-3)
}
