# go-many-git
Tool to manage multiple git repositories

## Installation
```
go get -u github.com/abrochard/go-many-git
go install github.com/abrochard/go-many-git
echo 'alias gmg=$GOPATH/bin/go-many-git' >> ~/.bashrc
```

## Usage
```
Usage: gmg [@tag] <comand> [<args>]

Go-many-git basic usage is to run a particular git command across multiple repos
For example, 'gmg pull' runs 'git pull' across all registered repos
By default 'gmg' alone runs 'git status'

Optionally, a repos can be identified by a shared tag (@example), making it possible to target a subset of repos

Go-many-git accepts all git commands, but here are a few gmg specific commands:

   [@tag] register <path>    Add the repo in <path> to the list of repos, with an optional tag
   unregister <path>         Remove the repo in <path> from the list
   list                      Print all registered repos
   help                      Print this help
```
