# SSH Connection Manager

A CLI tool to manage SSH connections and connect to remote servers using aliases.

## Build

```bash
go build -o ./dist/scm ./cmd/main.go

```

OR

```bash
docker run --rm -v "${PWD}:/usr/src/app" -w /usr/src/app golang:1.22 \
    env GOOS=linux GOARCH=amd64 go build -o dist/scm ./cmd/main.go

docker run --rm -v "${PWD}:/usr/src/app" -w /usr/src/app golang:1.22 \
    env GOOS=windows GOARCH=amd64 go build -o dist\scm .\cmd\main.go

docker run --rm -v "${PWD}:/usr/src/app" -w /usr/src/app golang:1.22 \
    env GOOS=darwin GOARCH=amd64 go build -o dist/scm ./cmd/main.go

```
