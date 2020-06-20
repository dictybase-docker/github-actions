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

## help, h

Shows a list of commands or help for one command
 <nil>
