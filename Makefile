run:
	go run ./cmd/tui

deps:
	go mod tidy
	go mod download
