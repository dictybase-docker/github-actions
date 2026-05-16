# github-actions

[![License](https://img.shields.io/badge/License-BSD%202--Clause-blue.svg)](LICENSE)
![Continuous integration](https://github.com/dictybase-docker/github-actions/actions/workflows/ci.yaml/badge.svg?branch=develop)
[![Code coverage](https://codecov.io/gh/dictybase-docker/github-actions/branch/develop/graph/badge.svg)](https://codecov.io/gh/dictybase-docker/github-actions)
![Last commit](https://badgen.net/github/last-commit/dictybase-docker/github-actions/develop)
[![Funding](https://badgen.net/badge/Funding/Rex%20L%20Chisholm,dictyBase,DCR/yellow?list=|)](https://reporter.nih.gov/project-details/10024726)

![Go](https://badgen.net/static/go/1.25/00ADD8)
[![GitHub release](https://badgen.net/github/release/dictybase-docker/github-actions)](https://github.com/dictybase-docker/github-actions/releases)

CLI to manage dictyBase GitHub actions, deployments, and repository operations.

## Contents

- [Available commands](#available-commands)
- [Commands](#commands)
  - [issue-comment-report](#issue-comment-report)
  - [issue-comment-count](#issue-comment-count)
  - [store-report](#store-report)
  - [deploy-status](#deploy-status)
  - [share-deploy-payload](#share-deploy-payload)
  - [get-cluster-credentials](#get-cluster-credentials)
  - [deploy-chart](#deploy-chart)
  - [files-committed](#files-committed)
  - [batch-multi-repo](#batch-multi-repo)
  - [parse-chatops-deploy](#parse-chatops-deploy)
  - [report-as-comment](#report-as-comment)
  - [migrate-repos](#migrate-repos)
  - [analytics-report](#analytics-report)
  - [setup-dagger-checksum](#setup-dagger-checksum)
  - [setup-dagger-bin](#setup-dagger-bin)
- [Full CLI documentation](#full-cli-documentation)

## Available commands

```
NAME:
   gh-action - run github action

USAGE:
   github-actions [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
   issue-comment-report, icr      reports no of comments for every issue
   issue-comment-count, icc       reports total no of issues and comments since a particular date
   store-report, ur               save report to s3 storage
   deploy-status, ds              create a github deployment status
   share-deploy-payload, sdp      share deployment payload data in github workflow
   get-cluster-credentials, gcre  get kubernetes cluster credentials using gcloud
   doc                            generate markdown documentation
   deploy-chart, dc               deploy helm chart
   files-committed, fc            outputs list of file committed in a git push or pull-request
   batch-multi-repo, bmr          Commit a file to multiple repositories
   parse-chatops-deploy, pcd      parses chatops deploy command and extracts ref and image tag values
   report-as-comment, rac         generate ontology report in pull request comment
   migrate-repos, mr              fork and migrate repositories to a different owner or organization
   analytics-report, ar           generate google analytics report of sessions,users,pageviews in csv format.
   setup-dagger-checksum, sc      setup checksum of dagger binary
   setup-dagger-bin, sd           setup dagger command line
   help, h                        Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --log-format value            format of the log, either of json or text. (default: "json")
   --log-level value             log level for the application (default: "error")
   --token value, -t value       github personal access token [$GITHUB_TOKEN]
   --repository value, -r value  Github repository
   --owner value                 Github repository owner (default: "dictyBase")
   --help, -h                    show help
   --version, -v                 print the version
```

## Commands

### issue-comment-report

Reports number of comments for every issue in a repository. Outputs a CSV file.

```
gh-action issue-comment-report [--token TOKEN] [--repository REPO] [--owner OWNER] [--output FILE] [--state STATE]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--output` | *(auto-generated)* | CSV output file path |
| `--state` | `all` | Issue state filter (`open`, `closed`, `all`) |

### issue-comment-count

Reports total number of issues and comments since a particular date. Outputs a CSV file.

```
gh-action issue-comment-count --date DATE [--output FILE]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--date, -d` | *(required)* | Start date for counting (e.g. `2024-01-01`) |
| `--output, -o` | `stdout` | Output file path |

### store-report

Upload a report file to S3-compatible storage.

```
gh-action store-report --input FILE --upload-path PATH [--access-key KEY] [--secret-key KEY] [--s3-server HOST] [--s3-server-port PORT] [--s3-bucket BUCKET]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--input` | *(required)* | Input file to upload |
| `--upload-path, -p` | *(required)* | Full upload path inside the bucket |
| `--access-key, --akey` | *(auto)* | S3 access key |
| `--secret-key, --skey` | *(auto)* | S3 secret key |
| `--s3-server` | `minio` | S3 server endpoint |
| `--s3-server-port` | *(auto)* | S3 server port |
| `--s3-bucket` | `report` | S3 bucket name |

### deploy-status

Create a GitHub deployment status for a given deployment.

```
gh-action deploy-status --deployment_id ID --state STATE [--url URL]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--deployment_id` | *(required)* | Deployment identifier |
| `--state` | *(required)* | Deployment state |
| `--url` | `""` | URL associated with this status |

### share-deploy-payload

Share deployment payload data in a GitHub workflow by writing key-value pairs to `$GITHUB_OUTPUT`.

```
gh-action share-deploy-payload --payload-file FILE
```

| Flag | Default | Description |
|------|---------|-------------|
| `--payload-file, -f` | *(required)* | Path to the JSON deploy payload file |

### get-cluster-credentials

Get Kubernetes cluster credentials using gcloud.

```
gh-action get-cluster-credentials --cluster NAME --project ID --zone ZONE
```

| Flag | Default | Description |
|------|---------|-------------|
| `--cluster` | *(required)* | Kubernetes cluster name |
| `--project` | *(required)* | Google Cloud project ID |
| `--zone` | *(required)* | Compute zone for the cluster |

### deploy-chart

Deploy a Helm chart to a Kubernetes cluster.

```
gh-action deploy-chart --name NAME --path PATH --namespace NS --image-tag TAG
```

| Flag | Default | Description |
|------|---------|-------------|
| `--name` | *(required)* | Chart name |
| `--path` | *(required)* | Relative chart path from repo root |
| `--namespace` | *(required)* | Kubernetes namespace |
| `--image-tag` | *(required)* | Docker image tag |

### files-committed

Output list of files committed in a git push or pull request event.

```
gh-action files-committed --payload-file FILE [--event-type TYPE] [--output FILE] [--include-file-suffix SUFFIX] [--skip-deleted]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--payload-file, -f` | *(required)* | Path to the event payload JSON |
| `--event-type, -e` | `push` | Event type (`push` or `pull-request`) |
| `--output, -o` | `stdout` | Output file path |
| `--include-file-suffix, --ifs` | `""` | Only report files with this suffix |
| `--skip-deleted, --sd` | `false` | Skip deleted files |

### batch-multi-repo

Commit a file to multiple repositories in a single operation.

```
gh-action batch-multi-repo --input-file FILE --repository-list FILE --repository-path PATH [--branch BRANCH]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--input-file, -i` | *(required)* | File to commit |
| `--repository-list, -l` | *(required)* | File with list of repository names (one per line) |
| `--repository-path, --rp` | *(required)* | Relative path in target repositories |
| `--branch, -b` | `develop` | Target branch (must already exist) |

### parse-chatops-deploy

Parse a ChatOps deploy command and extract ref and image tag values.

```
gh-action parse-chatops-deploy --payload-file FILE [--frontend]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--payload-file, -f` | *(required)* | Path to JSON payload |
| `--frontend` | `false` | Set when deploying a frontend web app |

### report-as-comment

Generate an ontology report and post it as a pull request comment.

```
gh-action report-as-comment --commit-list-file FILE --pull-request-id ID --report-dir DIR [--ref REF]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--commit-list-file, -c` | *(required)* | File with list of committed files |
| `--pull-request-id, --id` | *(required)* | Pull request ID to comment on |
| `--report-dir, -d` | *(required)* | Directory containing ontology reports |
| `--ref` | `""` | Git reference (tag or commit ID) |

### migrate-repos

Fork and migrate repositories to a different owner or organization.

```
gh-action migrate-repos --owner-to-migrate OWNER --repo-to-move REPO [--repo-to-move REPO...] [--poll-for SECONDS] [--poll-interval SECONDS]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--owner-to-migrate, --om` | *(required)* | Target owner/organization name |
| `--repo-to-move, -m` | *(required, repeatable)* | Repository names to migrate |
| `--poll-for` | `60` | Max polling duration in seconds |
| `--poll-interval` | `2` | Polling interval in seconds |

### analytics-report

Generate a Google Analytics report (sessions, users, pageviews) in CSV format.

```
gh-action analytics-report --view-id ID --credential-file FILE --start-date DATE [--end-date DATE] [--output FILE]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--view-id` | *(required)* | Google Analytics view ID |
| `--credential-file, -c` | *(required)* | Google service account credential file |
| `--start-date, -s` | *(required)* | Start date (`YYYY-MM-DD` format) |
| `--end-date, -e` | `today` | End date (`YYYY-MM-DD` format) |
| `--output, -o` | `stdout` | Output CSV file name |

### setup-dagger-checksum

Verify the Dagger binary checksum from stdin.

```
gh-action setup-dagger-checksum < checksum_file
```

### setup-dagger-bin

Set up the Dagger CLI binary and add it to `$PATH`.

```
gh-action setup-dagger-bin
```

## Full CLI documentation

Run `gh-action doc` to generate the complete CLI reference in markdown format with all flags and options.
