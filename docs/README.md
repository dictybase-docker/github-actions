% gh-action(8) 

% 

# NAME

gh-action - run github action

# SYNOPSIS

gh-action

```
[--help|-h]
[--log-format]=[value]
[--log-level]=[value]
[--owner]=[value]
[--repository|-r]=[value]
[--token|-t]=[value]
[--version|-v]
```

**Usage**:

```
gh-action [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--help, -h**: show help

**--log-format**="": format of the log, either of json or text. (default: json)

**--log-level**="": log level for the application (default: error)

**--owner**="": Github repository owner (default: dictyBase)

**--repository, -r**="": Github repository

**--token, -t**="": github personal access token

**--version, -v**: print the version


# COMMANDS

## issue-comment-report, icr

reports no of comments for every issue

**--output**="": file where csv format output is written, creates a timestamp based file by default

**--state**="": state of the issue for filtering (default: all)

## store-report, ur

save report to s3 storage

**--access-key, --akey**="": access key for S3 server, required based on command run

**--input**="": input file that will be uploaded

**--s3-bucket**="": S3 bucket where the data will be uploaded (default: report)

**--s3-server**="": S3 server endpoint (default: minio)

**--s3-server-port**="": S3 server port

**--secret-key, --skey**="": secret key for S3 server, required based on command run

**--upload-path, -p**="": full upload path inside the bucket

## deploy-status, ds

create a github deployment status

**--deployment_id**="": Deployment identifier (default: 0)

**--state**="": The state of the deployment status

**--url**="": The url that is associated with this status

## share-deploy-payload, sdp

share deployment payload data in github workflow

**--payload-file, -f**="": Full path to the file that contain the deploy payload

## get-cluster-credentials, gcre

get kubernetes cluster credentials using gcloud

**--cluster**="": Name of k8s cluster

**--project**="": Google cloud project id

**--zone**="": Compute zone for the cluster

## doc

generate markdown documentation

## deploy-chart, dc

deploy helm chart

**--image-tag**="": Docker image tag

**--name**="": Name of the chart

**--namespace**="": Kubernetes namespace

**--path**="": Relative chart path from the root of the repo

## files-committed, fc

outputs list of file committed in a git push or pull-request

**--event-type, -e**="": Name of the event, either or push or pull-request (default: push)

**--include-file-suffix, --ifs**="": file with the given suffix will only be reported

**--output, -o**="": Name of output file, defaults to stdout

**--payload-file, -f**="": Full path to the file that contain the event payload

**--skip-deleted, --sd**: skip deleted files in the commit(s)

## batch-multi-repo, bmr

Commit a file to multiple repositories

**--branch, -b**="": repository branch(should exist before committing) (default: develop)

**--input-file, -i**="": file that will be committed to repository

**--repository-list, -l**="": file with list of repositories name, one repository in every line

**--repository-path, --rp**="": relative path(from root) in the repository for the input file

## parse-chatops-deploy, pcd

parses chatops deploy command and extracts ref and image tag values

**--frontend**: used if deploying frontend web app (needed for updating image-tag correctly)

**--payload-file, -f**="": path to JSON payload

## report-as-comment, rac

generate ontology report in pull request comment

**--commit-list-file, -c**="": file that contain a list of committed file in the push event

**--pull-request-id, --id**="": id of a pull-request where the comment should be made (default: 0)

**--report-dir, -d**="": folder containing ontology reports

## help, h

Shows a list of commands or help for one command
 <nil>
% gh-action(8) 

% 

# NAME

gh-action - run github action

# SYNOPSIS

gh-action

```
[--help|-h]
[--log-format]=[value]
[--log-level]=[value]
[--owner]=[value]
[--repository|-r]=[value]
[--token|-t]=[value]
[--version|-v]
```

**Usage**:

```
gh-action [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--help, -h**: show help

**--log-format**="": format of the log, either of json or text. (default: json)

**--log-level**="": log level for the application (default: error)

**--owner**="": Github repository owner (default: dictyBase)

**--repository, -r**="": Github repository

**--token, -t**="": github personal access token

**--version, -v**: print the version


# COMMANDS

## issue-comment-report, icr

reports no of comments for every issue

**--output**="": file where csv format output is written, creates a timestamp based file by default

**--state**="": state of the issue for filtering (default: all)

## store-report, ur

save report to s3 storage

**--access-key, --akey**="": access key for S3 server, required based on command run

**--input**="": input file that will be uploaded

**--s3-bucket**="": S3 bucket where the data will be uploaded (default: report)

**--s3-server**="": S3 server endpoint (default: minio)

**--s3-server-port**="": S3 server port

**--secret-key, --skey**="": secret key for S3 server, required based on command run

**--upload-path, -p**="": full upload path inside the bucket

## deploy-status, ds

create a github deployment status

**--deployment_id**="": Deployment identifier (default: 0)

**--state**="": The state of the deployment status

**--url**="": The url that is associated with this status

## share-deploy-payload, sdp

share deployment payload data in github workflow

**--payload-file, -f**="": Full path to the file that contain the deploy payload

## get-cluster-credentials, gcre

get kubernetes cluster credentials using gcloud

**--cluster**="": Name of k8s cluster

**--project**="": Google cloud project id

**--zone**="": Compute zone for the cluster

## doc

generate markdown documentation

## deploy-chart, dc

deploy helm chart

**--image-tag**="": Docker image tag

**--name**="": Name of the chart

**--namespace**="": Kubernetes namespace

**--path**="": Relative chart path from the root of the repo

## files-committed, fc

outputs list of file committed in a git push or pull-request

**--event-type, -e**="": Name of the event, either or push or pull-request (default: push)

**--include-file-suffix, --ifs**="": file with the given suffix will only be reported

**--output, -o**="": Name of output file, defaults to stdout

**--payload-file, -f**="": Full path to the file that contain the event payload

**--skip-deleted, --sd**: skip deleted files in the commit(s)

## batch-multi-repo, bmr

Commit a file to multiple repositories

**--branch, -b**="": repository branch(should exist before committing) (default: develop)

**--input-file, -i**="": file that will be committed to repository

**--repository-list, -l**="": file with list of repositories name, one repository in every line

**--repository-path, --rp**="": relative path(from root) in the repository for the input file

## parse-chatops-deploy, pcd

parses chatops deploy command and extracts ref and image tag values

**--frontend**: used if deploying frontend web app (needed for updating image-tag correctly)

**--payload-file, -f**="": path to JSON payload

## report-as-comment, rac

generate ontology report in pull request comment

**--commit-list-file, -c**="": file that contain a list of committed file in the push event

**--pull-request-id, --id**="": id of a pull-request where the comment should be made (default: 0)

**--ref**="": git reference name, tag or commit id

**--report-dir, -d**="": folder containing ontology reports

## migrate-repos, mr

fork and migrate repositories to a different owner or organization

**--owner-to-migrate, --om**="": owner name where the repositories will be migrated

**--poll-for**="": threshold for polling forked repository(in seconds) (default: 60)

**--poll-interval**="": polling interval for forked repository(in seconds) (default: 2)

**--repo-to-move, -m**="": repositories that will be migrated, repeat the option for multiple values

## help, h

Shows a list of commands or help for one command
 <nil>
