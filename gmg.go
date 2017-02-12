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

func printHelp() {
	fmt.Println("Usage")
}

func main() {

	args := os.Args[1:]
	tag := ""

	if len(args) == 0 {
		args = []string{"status"}
	} else if string(args[0][0]) == "@" {
		tag, args = args[0][1:], args[1:]
	}
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
		path, _ := filepath.Abs(args[1])
		registerRepo(path, tag, repos)
	case "unregister":
		path, _ := filepath.Abs(args[1])
		unregisterRepo(path, repos)
	case "help":
	case "--help":
	case "-h":
		printHelp()
	default:
		runCmd(repos, tag, args...)
	}

}
