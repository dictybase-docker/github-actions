FROM golang:1.11.13-alpine3.10
LABEL maintainer="Siddhartha Basu <siddhartha-basu@northwestern.edu>"
RUN apk add --no-cache git build-base
RUN mkdir -p /github-actions
WORKDIR /github-actions
COPY go.mod ./
COPY go.sum ./
RUN go mod download
ADD cmd cmd
ADD internal internal
RUN go build -o app cmd/import/main.go

FROM alpine:3.10
RUN apk --no-cache add ca-certificates
COPY --from=0 /github-content/app /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/app"]
