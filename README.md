# httpchecker

CLI tool in Go to check HTTP status codes of a list of URLs with concurrency, timeout, and retries.

## Installation via GitHub (go install)

After publishing to GitHub, install with:

```bash
go install github.com/dcaiovinicius/httpchecker/cmd/httpchecker@latest
```

### How to use

```
echo 'https://www.google.com' | httpcheker
```

```
httpcheker https://www.google.com
```