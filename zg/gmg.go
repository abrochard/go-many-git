package zg

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

// Location of the config file under the $HOME dir
var ConfigFile = "/.config/zg-repos.json"

// Repo struct
type Repo struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Tag      string `json:"tag"`
}

// Check for errors, print message and panic if needed
func check(err error, message string) {
	if err != nil {
		log.Fatalf("%v, error: %v", message, err)
	}
}

// Set the path to the config file
func SetConfigFilePath() {
	usr, err := user.Current()
	check(err, "Error getting the current user")
	ConfigFile = usr.HomeDir + ConfigFile
}

// Create the config file with empty repo list
func createConfigFile() {
	data := []byte("[]\n")
	err := ioutil.WriteFile(ConfigFile, data, 0644)
	check(err, "Failed to create config file")
}

// Deserialize the repos from the config file into structs
func GetRepos() []Repo {
	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
		// config file does not exits
		createConfigFile()
		return make([]Repo, 0)
	}

	raw, err := ioutil.ReadFile(ConfigFile)
	check(err, "Config file not found")

	var c []Repo
	_ = json.Unmarshal(raw, &c)
	return c
}

// Serialize the repos structs into the config file
func saveRepos(repos []Repo) {
	buff, err := json.MarshalIndent(repos, "", "   ")

	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(ConfigFile, buff, 0644)
	check(err, "Failed to save repos to config file")
}

// Run a git command across all repos with matching tag
func RunCmd(repos []Repo, tag string, args ...string) {
	for _, r := range repos {
		if tag == "" || (tag != "" && r.Tag != "" && tag == r.Tag) {
			var cmd *exec.Cmd

			if args[0] == "-exec" {
				cmd = exec.Command(args[1], args[2:]...)
			} else {
				params := []string{"-C", r.Location}
				params = append(params, args...)
				cmd = exec.Command("git", params...)
			}

			var out bytes.Buffer
			var stderr bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &stderr

			color.Cyan(r.Name)
			err := cmd.Run()

			if err != nil {
				fmt.Println(stderr.String())
			}

			fmt.Printf("%s\n", out.String())
		}
	}
}

// Add a new repo to the list
func RegisterRepo(location string, tag string, repos []Repo) {
	_, err := os.Stat(location)
	check(err, "Invalid path")
	name := filepath.Base(location)
	r := Repo{Name: name, Location: location, Tag: tag}
	repos = append(repos, r)
	saveRepos(repos)
}

// Remove a repo from the list
func UnregisterRepo(path string, repos []Repo) {
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
func PrintRepos(tag string, repos []Repo) {
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
func PrintHelp() {
	fmt.Println("Usage: ziti-git [@tag] <command> [<args>]")
	fmt.Println("")
	fmt.Println("ziti-git basic usage is to run a particular git command across multiple repos")
	fmt.Println("For example, 'ziti-git pull' runs 'git pull' across all registered repos")
	fmt.Println("By default 'ziti-git' alone runs 'git status'")
	fmt.Println("")
	fmt.Println("Optionally, a repos can be identified by a shared tag (@example), making it possible to target a subset of repos")
	fmt.Println("ie: `ziti-git @api pull` runs `git pull` on all repos tagged with `api`")
	fmt.Println("")
	fmt.Println("Arbitrary non-git commands can be run by using the -exec flag")
	fmt.Println("")
	fmt.Println("ziti-git accepts all git commands, but here are a few ziti-git specific commands:")
	fmt.Println("")
	fmt.Println("   [@tag] register <path>    Add the repo in <path> to the list of repos, with an optional tag")
	fmt.Println("   unregister <path>         Remove the repo in <path> from the list")
	fmt.Println("   [@tag] -table-status      Will output the name, branch, and status of each repo in a table")
	fmt.Println("   [@tag] b                  Shorthand to display the repos current branch")
	fmt.Println("   [@tag] -exec ls           Execute non-git command ls on each repo")
	fmt.Println("   list                      Print all registered repos")
	fmt.Println("   help                      Print this help")
	fmt.Println("")
	fmt.Println("See README.md for more details")
}

// Check for valid tag
func ValidTag(str string) bool {
	return string(str[0]) == "@" && len(str) > 1
}

func getTrailingEllipsesString(orig string, max int) string {
	if len(orig) > max {
		max -= 3
		return orig[0:len(orig)-max] + "..."
	}
	return orig
}

func getMiddleEllipsesString(orig string, max int) string {
	if len(orig) > max {
		max -= 3

		middle := len(orig) / 2
		remove := len(orig) - max

		numLeft := remove / 2
		numRight := remove - numLeft

		return orig[0:middle-numLeft] + "..." + orig[middle+numRight:]
	}
	return orig
}

func TableStatus(repos []Repo, tag string) {
	var statusData [][]string

	var errorData [][]string

	for _, r := range repos {
		if tag == "" || (tag != "" && r.Tag != "" && tag == r.Tag) {

			entry := []string{getMiddleEllipsesString(r.Name, 20), "", "", "", "", getMiddleEllipsesString(r.Location, 60)}

			branch, err := getCurrentBranch(r.Location)

			if err != nil {
				errNum := strconv.Itoa(len(errorData) + 1)
				errorData = append(errorData, []string{errNum, err.Error(), err.StdErr.String()})
				branch = "See Error: " + errNum
			}

			tag, _ := getCurrentTag(r.Location)

			if tag == "" {
				tag, _ = getCurrentCommit(r.Location)
			}

			tag = getTrailingEllipsesString(tag, 10)
			branch = getTrailingEllipsesString(branch, 25)

			entry[1] = branch
			entry[2] = tag

			stagedStatus, unstagedStatus, err := getStatus(r.Location)

			if err != nil {
				errNum := strconv.Itoa(len(errorData) + 1)
				errorData = append(errorData, []string{errNum, err.Error(), err.StdErr.String()})
				stagedStatus = "See Error: " + errNum
				unstagedStatus = "See Error: " + errNum
			}

			entry[3] = stagedStatus
			entry[4] = unstagedStatus

			statusData = append(statusData, entry)
		}
	}

	statusTable := tablewriter.NewWriter(os.Stdout)
	statusTable.SetHeader([]string{"Name", "Branch", "Tag", "Staged", "Unstaged", "Location"})
	statusTable.SetColumnAlignment([]int{tablewriter.ALIGN_DEFAULT, tablewriter.ALIGN_DEFAULT, tablewriter.ALIGN_DEFAULT, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_DEFAULT})
	statusTable.SetColumnColor(tablewriter.Colors{}, tablewriter.Colors{}, tablewriter.Colors{}, tablewriter.Colors{tablewriter.FgHiGreenColor}, tablewriter.Colors{tablewriter.FgHiRedColor}, tablewriter.Colors{})

	for _, data := range statusData {
		statusTable.Append(data)
	}
	statusTable.Render()

	if len(errorData) > 0 {
		errTable := tablewriter.NewWriter(os.Stdout)
		errTable.SetHeader([]string{"Error Number", "Message", "StdErr"})

		for _, data := range errorData {
			errTable.Append(data)
		}
		errTable.Render()
	}
}

type StdOutErr struct {
	StdErr bytes.Buffer
	Err    error
}

func (err *StdOutErr) Error() string {
	return err.Err.Error()
}

func getCurrentBranch(location string) (string, *StdOutErr) {
	params := []string{"-C", location, "rev-parse", "--abbrev-ref", "HEAD"}
	cmd := exec.Command("git", params...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		return "", &StdOutErr{
			StdErr: stderr,
			Err:    err,
		}
	}

	return strings.TrimSpace(out.String()), nil
}

func getCurrentTag(location string) (string, *StdOutErr) {
	params := []string{"-C", location, "describe", "--abbrev=0"}
	cmd := exec.Command("git", params...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		return "", &StdOutErr{
			StdErr: stderr,
			Err:    err,
		}
	}

	return strings.TrimSpace(out.String()), nil
}

func getCurrentCommit(location string) (string, *StdOutErr) {
	params := []string{"-C", location, "rev-parse", "--short", "HEAD"}
	cmd := exec.Command("git", params...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		return "", &StdOutErr{
			StdErr: stderr,
			Err:    err,
		}
	}

	return strings.TrimSpace(out.String()), nil
}

func getStatus(location string) (string, string, *StdOutErr) {
	params := []string{"-C", location, "status", "--porcelain=v2"}

	cmd := exec.Command("git", params...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		return "", "", &StdOutErr{
			StdErr: stderr,
			Err:    err,
		}
	}
	staged, notStaged := parseStatus(&out)

	return staged, notStaged, nil
}

func parseStatus(statusBuffer *bytes.Buffer) (string, string) {
	if statusBuffer.Len() == 0 {
		return "", ""
	}

	scanner := bufio.NewScanner(statusBuffer)

	stagedDeletes := 0
	notStagedDeletes := 0

	stagedModifies := 0
	notStagedModifies := 0

	stagedNew := 0
	notStagedNew := 0

	for scanner.Scan() {
		line := scanner.Text()
		values := strings.Split(line, " ")

		if values[0] == "?" {
			notStagedNew++
			break
		}

		switch values[1] {
		case "D.":
			stagedDeletes++
		case ".D":
			notStagedDeletes++
		case "M.":
			stagedModifies++
		case ".M":
			notStagedModifies++
		case "R.":
			stagedModifies++
		case ".R":
			notStagedModifies++
		case "A.":
			stagedNew++
		}
	}

	stagedStr := ""

	if stagedNew != 0 || stagedModifies != 0 || stagedDeletes != 0 {
		stagedStr = fmt.Sprintf("+%d ~%d -%d", stagedNew, stagedModifies, stagedDeletes)
	}

	notStagedStr := ""

	if notStagedNew != 0 || notStagedModifies != 0 || notStagedDeletes != 0 {
		notStagedStr = fmt.Sprintf("+%d ~%d -%d", notStagedNew, notStagedModifies, notStagedDeletes)
	}

	return stagedStr, notStagedStr

}
