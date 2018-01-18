# go-github-reviewer-stats

## Usage

```
✘╹◡╹✘ < ./go-github-reviewer-stats --help
Usage of ./go-github-reviewer-stats:
  -base-url string
        custom GitHub base URL if you use GitHub Enterprise (default "https://api.github.com")
  -insecure-skip-verify
        skip verification of cert
  -owner string
        owner of repo
  -per-page int
        count of pull requests to scan (default 10)
  -repo string
        repo name
```

## Build

```sh
dep ensure
go build ./...
```
