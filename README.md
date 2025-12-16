# dua

Disk Usage Analyzer

## Background

Inspired by [dua](https://github.com/Byron/dua-cli), this project explores a similar idea in Go, focusing on learning and personal experimentation rather than reproducing the original implementation.

---

## Local Environment

- Build and run the application

```sh
go build -o bin/dua main.go
./bin/dua scan .
```

Basic commands:

- Run linter

```sh
golangci-lint run
```

- Run test

```sh
go test -v -coverprofile=coverage.out ./...
```

To ensure consistent results between your local environment and the repository’s CI pipeline, use the same tool versions. For example:

```
- name: Run golangci-lint
  uses: golangci/golangci-lint-action@v8
  with:
    version: v2.4.0
```

The action [golangci-lint-action@v8](https://github.com/golangci/golangci-lint-action/tree/v8) defaults to [golangci-lint v2.1.0](https://github.com/golangci/golangci-lint/tree/v2.1.0), but the `version` field overrides it to use [v2.4.0](https://github.com/golangci/golangci-lint/tree/v2.4.0).