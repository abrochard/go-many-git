package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

// Location of the config file under the $HOME dir
var CONFIG_FILE = "/.config/gmg-repos.json"

// Repo struct
type Repo struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Tag      string `json:"tag"`
}

// Check for errors, print message and panic if needed
func check(err error, message string) {
	if err != nil {
		fmt.Println(message)
		panic(err)
	}
}

// Set the path to the config file
func setConfigFilePath() {
	usr, err := user.Current()
	check(err, "Error getting the current user")
	CONFIG_FILE = usr.HomeDir + CONFIG_FILE
}

// Create the config file with empty repo list
func createConfigFile() {
	data := []byte("[]\n")
	err := ioutil.WriteFile(CONFIG_FILE, data, 0644)
	check(err, "Failed to create config file")
}

// Deserialize the repos from the config file into structs
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

// Serialize the repos structs into the config file
func saveRepos(repos []Repo) {
	bytes, err := json.MarshalIndent(repos, "", "   ")

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(CONFIG_FILE, bytes, 0644)
	check(err, "Failed to save repos to config file")
}

// Run a git command across all repos with matching tag
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

// Add a new repo to the list
func registerRepo(location string, tag string, repos []Repo) {
	name := filepath.Base(location)
	r := Repo{Name: name, Location: location, Tag: tag}
	repos = append(repos, r)
	saveRepos(repos)
}

// Remove a repo from the list
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

// Pretty print all repos in list
func printRepos(tag string, repos []Repo) {
	for _, r := range repos {
		if tag == "" || (tag != "" && r.Tag != "" && tag == r.Tag) {
			fmt.Printf("Name: ")
			color.Cyan("%s", r.Name)
			fmt.Printf("Location: %s\n", r.Location)
			fmt.Printf("Tag: %s\n", r.Tag)
			fmt.Println("")
		}
	}
}

// Print help message
func printHelp() {
	fmt.Println("Usage: gmg [@tag] <comand> [<args>]")
	fmt.Println("")
	fmt.Println("Go-many-git basic usage is to run a particular git command across multiple repos")
	fmt.Println("For example, 'gmg pull' runs 'git pull' across all registered repos")
	fmt.Println("By default 'gmg' alone runs 'git status'")
	fmt.Println("")
	fmt.Println("Optionally, a repos can be identified by a shared tag (@example), making it possible to target a subset of repos")
	fmt.Println("ie: `gmg @api pull` runs `git pull` on all repos tagged with `api`")
	fmt.Println("")
	fmt.Println("Go-many-git accepts all git commands, but here are a few gmg specific commands:")
	fmt.Println("")
	fmt.Println("   [@tag] register <path>    Add the repo in <path> to the list of repos, with an optional tag")
	fmt.Println("   unregister <path>         Remove the repo in <path> from the list")
	fmt.Println("   [@tag] b                  Shorthand to display the repos current branch")
	fmt.Println("   list                      Print all registered repos")
	fmt.Println("   help                      Print this help")
	fmt.Println("")
	fmt.Println("See README.md for more details")
}

// Check for valid tag
func validTag(str string) bool {
	return string(str[0]) == "@" && len(str) > 1
}

// Parse args into tag and commands
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

// Main
func main() {
	setConfigFilePath()

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
		// List repos
		printRepos(tag, repos)
	case "b":
		// List branches short
		args = []string{"rev-parse", "--abbrev-ref", "HEAD"}
		runCmd(repos, tag, args...)
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
