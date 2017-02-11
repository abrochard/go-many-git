package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Repo struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Tag      string `json:"tag"`
}

func (p Repo) toString() string {
	return toJson(p)
}

func toJson(p interface{}) string {
	bytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return string(bytes)
}

func getRepos() []Repo {
	raw, err := ioutil.ReadFile("./repos.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c []Repo
	json.Unmarshal(raw, &c)
	return c
}

func main() {

	repos := getRepos()

	args := os.Args[1:]
	var tag, action, cmd = "", "", ""
	l := len(args)
	if l == 0 {
		action = "status"
	} else if l == 1 {
		action = args[0]
	} else if l == 2 {
		if string(args[0][0]) == "@" {
			tag = args[0]
			action = args[1]
		} else {
			action = args[0]
			cmd = args[1]
		}
	} else if l >= 3 {
		tag = args[0]
		action = args[1]
		cmd = strings.Join(args[2:], " ")
	}

	fmt.Println(cmd)

	for _, r := range repos {
		if tag == "" || (tag != "" && r.Tag != "" && tag == r.Tag) {
			color.Cyan(r.Name)
			out, err := exec.Command("git", "-C", r.Location, action).Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("The output is %s\n", out)
		}
	}

}
