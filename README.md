# github-actions
[![License](https://img.shields.io/badge/License-BSD%202--Clause-blue.svg)](LICENSE)   
[![Technical debt](https://badgen.net/codeclimate/tech-debt/dictyBase-docker/github-actions)](https://codeclimate.com/github/dictyBase-docker/github-actions/trends/technical_debt)
[![Issues](https://badgen.net/codeclimate/issues/dictyBase-docker/github-actions)](https://codeclimate.com/github/dictyBase-docker/github-actions/issues)
[![Maintainability](https://api.codeclimate.com/v1/badges/27d8dea5aa1373847404/maintainability)](https://codeclimate.com/github/dictybase-docker/github-actions/maintainability)
[![Dependabot Status](https://api.dependabot.com/badges/status?host=github&repo=dictyBase-docker/github-actions)](https://dependabot.com)   
![Issues](https://badgen.net/github/issues/dictyBase-docker/github-actions)
![Open Issues](https://badgen.net/github/open-issues/dictyBase-docker/github-actions)
![Closed Issues](https://badgen.net/github/closed-issues/dictyBase-docker/github-actions)   
![Total PRS](https://badgen.net/github/prs/dictyBase-docker/github-actions)
![Open PRS](https://badgen.net/github/open-prs/dictyBase-docker/github-actions)
![Closed PRS](https://badgen.net/github/closed-prs/dictyBase-docker/github-actions)
![Merged PRS](https://badgen.net/github/merged-prs/dictyBase-docker/github-actions)   
![Commits](https://badgen.net/github/commits/dictyBase-docker/github-actions/develop)
![Last commit](https://badgen.net/github/last-commit/dictyBase-docker/github-actions/develop)
![Branches](https://badgen.net/github/branches/dictyBase-docker/github-actions)
![Tags](https://badgen.net/github/tags/dictyBase-docker/github-actions/?color=cyan)   
![GitHub repo size](https://img.shields.io/github/repo-size/dictyBase-docker/github-actions?style=plastic)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/dictyBase-docker/github-actions?style=plastic)
[![Lines of Code](https://badgen.net/codeclimate/loc/dictyBase-docker/github-actions)](https://codeclimate.com/github/dictyBase-docker/github-actions/code)   
[![Funding](https://badgen.net/badge/NIGMS/Rex%20L%20Chisholm,dictyBase-docker/yellow?list=|)](https://projectreporter.nih.gov/project_info_description.cfm?aid=9476993)
[![Funding](https://badgen.net/badge/NIGMS/Rex%20L%20Chisholm,DSC/yellow?list=|)](https://projectreporter.nih.gov/project_info_description.cfm?aid=9438930)

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
