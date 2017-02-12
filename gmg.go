package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const CONFIG_FILE = "./repos.json"

type Repo struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Tag      string `json:"tag"`
}

func check(err error, message string) {
	if err != nil {
		fmt.Println(message)
		panic(err)
	}
}

func createConfigFile() {
	data := []byte("[]\n")
	err := ioutil.WriteFile(CONFIG_FILE, data, 0644)
	check(err, "Failed to create config file")
}

func getRepos() []Repo {
	if _, err := os.Stat(CONFIG_FILE); os.IsNotExist(err) {
		// config file does not exits
		createConfigFile()
		return make([]Repo, 0)
	}

	raw, err := ioutil.ReadFile(CONFIG_FILE)
	check(err, "Config file not found")

	var c []Repo
	json.Unmarshal(raw, &c)
	return c
}

func saveRepos(repos []Repo) {
	bytes, err := json.MarshalIndent(repos, "", "   ")

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(CONFIG_FILE, bytes, 0644)
	check(err, "Failed to save repos to config file")
}

func runCmd(repos []Repo, tag string, args ...string) {
	for _, r := range repos {
		if tag == "" || (tag != "" && r.Tag != "" && tag == r.Tag) {
			params := []string{"-C", r.Location}
			params = append(params, args...)
			color.Cyan(r.Name)
			out, err := exec.Command("git", params...).Output()
			check(err, "Failed to execute command")
			fmt.Printf("%s\n", out)
		}
	}
}

func registerRepo(location string, tag string, repos []Repo) {
	name := filepath.Base(location)
	r := Repo{Name: name, Location: location, Tag: tag}
	repos = append(repos, r)
	saveRepos(repos)
}

func unregisterRepo(path string, repos []Repo) {
	index := -1
	for i, r := range repos {
		if r.Location == path {
			index = i
			break
		}
	}

	if index >= 0 {
		repos = append(repos[:index], repos[index+1:]...)
		saveRepos(repos)
	}
}

func printRepos(tag string, repos []Repo) {
	for _, r := range repos {
		if tag == "" || (tag != "" && r.Tag != "" && tag == r.Tag) {
			fmt.Printf("Name: %s\n", r.Name)
			fmt.Printf("Location: %s\n", r.Location)
			fmt.Printf("Tag: %s\n", r.Tag)
			fmt.Println("")
		}
	}
}

func printHelp() {
	fmt.Println("Usage: gmg [@tag] <comand> [<args>]")
	fmt.Println("")
	fmt.Println("Go-many-git basic usage is to run a particular git command across multiple repos")
	fmt.Println("For example, 'gmg pull' runs 'git pull' across all registered repos")
	fmt.Println("By default 'gmg' alone runs 'git status'")
	fmt.Println("")
	fmt.Println("Optionally, a repos can be identified by a shared tag (@example), making it possible to target a subset of repos")
	fmt.Println("")
	fmt.Println("Go-many-git accepts all git commands, but here are a few gmg specific commands:")
	fmt.Println("")
	fmt.Println("   register <path> [@tag]    Add the repo in <path> to the list of repos, with an optional tag")
	fmt.Println("   unregister <path>         Remove the repo in <path> from the list")
	fmt.Println("   list                      Print all registered repos")
	fmt.Println("   help                      Print this help")
	fmt.Println("")
	fmt.Println("See README.md for more details")
}

func validTag(str string) bool {
	return string(str[0]) == "@" && len(str) > 1
}

func parseArgs(args []string) (string, []string) {
	tag := ""

	if len(args) == 0 {
		args = []string{"status"}
	} else if validTag(args[0]) {
		tag = args[0][1:] // scrub off the "@"
		if len(args) == 1 {
			args = []string{"status"}
		} else {
			args = args[1:]
		}
	}

	return tag, args
}

func main() {
	tag, args := parseArgs(os.Args[1:])
	cmd := args[0]

	repos := getRepos()
	if len(repos) == 0 && cmd != "register" {
		fmt.Println("No repositories registered. Nothing to do.")
		fmt.Println("Please register a repository with the command:")
		fmt.Println("gmg register [path]")
		os.Exit(0)
	}

	switch cmd {
	case "register":
		// Add this repo to the list
		path, _ := filepath.Abs(args[1])
		registerRepo(path, tag, repos)
	case "unregister":
		// Remove this repo from the list
		path, _ := filepath.Abs(args[1])
		unregisterRepo(path, repos)
	case "list":
		printRepos(tag, repos)
	case "help":
		// Print help
		printHelp()
	case "--help":
		// Print help
		printHelp()
	case "-h":
		// Print help
		printHelp()
	default:
		// Generic git command
		runCmd(repos, tag, args...)
	}

}
