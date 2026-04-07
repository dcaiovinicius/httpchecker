# httpchecker

CLI tool in Go to check HTTP status codes of a list of URLs with concurrency, timeout, and retries.

## Installation via GitHub (go install)

After publishing to GitHub, install with:

```bash
go install github.com/dcaiovinicius/httpchecker/cmd/httpchecker@v1.0.2
```

### How to use

```
echo 'https://www.google.com' | httpchecker
```

```
httpchecker https://www.google.com
```
