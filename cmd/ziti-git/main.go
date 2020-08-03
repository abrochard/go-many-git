package main

import (
	"fmt"
	"github.com/andrewpmartinez/ziti-git/zg"
	"os"
	"path/filepath"
)

func main() {

	zg.SetConfigFilePath()

	tag, args := parseArgs(os.Args[1:])
	cmd := args[0]

	repos := zg.GetRepos()
	if len(repos) == 0 && cmd != "register" {
		fmt.Println("No repositories registered. Nothing to do.")
		fmt.Println("Please register a repository with the command:")
		fmt.Println("ziti-git register [path]")
		os.Exit(0)
	}

	switch cmd {
	case "register":
		// Add this repo to the list
		if len(args) == 1 {
			fmt.Println("No path supplied")
			fmt.Println("ziti-git register [path]")
			os.Exit(1)
		}
		path, _ := filepath.Abs(args[1])
		zg.RegisterRepo(path, tag, repos)
	case "table-status":
		zg.TableStatus(repos, tag)
	case "ts":
		zg.TableStatus(repos, tag)
	case "unregister":
		// Remove this repo from the list
		path, _ := filepath.Abs(args[1])
		zg.UnregisterRepo(path, repos)
	case "list":
		// List repos
		zg.PrintRepos(tag, repos)
	case "b":
		// List branches short
		args = []string{"rev-parse", "--abbrev-ref", "HEAD"}
		zg.RunCmd(repos, tag, args...)
	case "help":
		// Print help
		zg.PrintHelp()
	case "--help":
		// Print help
		zg.PrintHelp()
	case "-h":
		// Print help
		zg.PrintHelp()
	default:
		// Generic git command
		zg.RunCmd(repos, tag, args...)
	}

}

// Parse args into tag and commands
func parseArgs(args []string) (string, []string) {
	tag := ""

	if len(args) == 0 {
		args = []string{"status"}
	} else if zg.ValidTag(args[0]) {
		tag = args[0][1:] // scrub off the "@"
		if len(args) == 1 {
			args = []string{"status"}
		} else {
			args = args[1:]
		}
	}

	return tag, args
}
