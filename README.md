# github-actions
cli to manage github actions 


## Commands

```
NAME:
   gh-action - run github action

USAGE:
   gact [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
     issue-comment-report, icr  reports no of comments for every issue
     help, h                    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --token value, -t value       github personal access token [$GITHUB_TOKEN]
   --repository value, -r value  Github repository
   --owner value                 Github repository owner (default: "dictyBase")
   --help, -h                    show help
   --version, -v                 print the version

```   

### Subcommands
```
NAME:
   gact issue-comment-report - reports no of comments for every issue

USAGE:
   gact issue-comment-report [command options] [arguments...]

OPTIONS:
   --output value  report output, goes to stdout by default
   --state value   state of the issue for filtering (default: "all")
```
