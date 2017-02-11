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

func createConfigFile() {
	data := []byte("[]\n")
	err := ioutil.WriteFile(CONFIG_FILE, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func getRepos() []Repo {
	if _, err := os.Stat(CONFIG_FILE); os.IsNotExist(err) {
		// config file does not exits
		createConfigFile()
		return make([]Repo, 0)
	}

	raw, err := ioutil.ReadFile(CONFIG_FILE)
	if err != nil {
		fmt.Println("Config file not found")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c []Repo
	json.Unmarshal(raw, &c)
	return c
}

func saveRepos(repos []Repo) {
	bytes, err := json.Marshal(repos)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(CONFIG_FILE, bytes, 0644)

	if err != nil {
		panic(err)
	}
}

func runCmd(repos []Repo, tag string, args ...string) {
	for _, r := range repos {
		if tag == "" || (tag != "" && r.Tag != "" && tag == r.Tag) {
			params := []string{"-C", r.Location}
			params = append(params, args...)
			color.Cyan(r.Name)
			out, err := exec.Command("git", params...).Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s\n", out)
		}
	}
}

func registerRepo(location string, repos []Repo) {

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

	repos := getRepos()
	if len(repos) == 0 {
		fmt.Println("No repositories registered. Nothing to do.")
		fmt.Println("Please register a repository with the command:")
		fmt.Println("gmg register [path]")
		os.Exit(0)
	}
	cmd := args[0]

	switch cmd {
	case "status":
		runCmd(repos, tag, args...)
	case "pull":
		runCmd(repos, tag, args...)
	case "push":
		runCmd(repos, tag, args...)
	case "cmd":
		args = args[1:] // Shift out the "cmd"
		runCmd(repos, tag, args...)
	case "register":
		path, _ := filepath.Abs(args[1])
		registerRepo(path, repos)
	case "unregister":
		path, _ := filepath.Abs(args[1])
		unregisterRepo(path, repos)
	case "help":
		printHelp()
	default:
		printHelp()
	}

}
