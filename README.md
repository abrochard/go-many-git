# ziti-git
A tool to manage multiple repositories with special considerations for the github.com/openziti project

## Requirements
- git >= 1.14

## Installation
```
go get -u github.com/andrewpmartinez/ziti-git
go install github.com/andrewpmartinez/ziti-git
```

## Alias to `zg`
```
echo 'alias zg=$GOPATH/bin/ziti-git' >> ~/.bashrc
```

## Usage
```
Usage: ziti-git [@tag] <comand> [<args>]

`ziti-git` basic usage is to run a particular git command across multiple repos
For example, 'ziti-git pull' runs 'git pull' across all registered repos
By default 'ziti-git' alone runs 'git status'

Additionally, arbitrary non-git commands can be run by specifying the `-exec` flag.

Optionally, a repos can be identified by a shared tag (@example), making it possible to target a subset of repos
ie: `ziti-git @api pull` runs `git pull` on all repos tagged with `api`

ziti-git accepts all git commands, but here are a few ziti-git specific commands:

   [@tag] register <path>    Add the repo in <path> to the list of repos, with an optional tag
   unregister <path>         Remove the repo in <path> from the list
   [@tag] table-status       Will output the name, branch, and status of each repo in a table
   [@tag] ts                 Will output the name, branch, and status of each repo in a table
   [@tag] b                  Shorthand to display the repos current branch
   [@tag] -exec ls           Run the non-git command ls on each repo
   list                      Print all registered repos
   help                      Print this help
```
Note that ziti-git writes its repo list to a config file located under `~/.config/zg-repos.json`

## Acknowledgements
Ziti Git is based off of [gmg](https://github.com/abrochard/go-many-git) which in turn was inspired by the amazing [mr](https://myrepos.branchable.com) and [gr](https://github.com/mixu/gr) tools.

A big thanks to all.
