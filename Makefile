build:
	go build -o bin/httpchecker ./cmd/httpchecker/main.go

run:
	go run ./cmd/httpchecker/main.go

linux:
	GOOS=linux GOARCH=amd64 go build -o bin/httpchecker ./cmd/httpchecker/main.go