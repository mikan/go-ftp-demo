.DEFAULT_GOAL := help
MAIN_PATH := cmd/tiny-ftp/*.go
OUT_PREFIX := build/tiny-ftp

.PHONY: clean
clean: ## Remove build artifact directory
	-rm -rfv build

.PHONY: lint
lint: ## Runs static code analysis
	go vet ./...
	command -v golint >/dev/null 2>&1 || { go get -u golang.org/x/lint/golint; }
	golint -set_exit_status ./...

.PHONY: run
run: ## Run tiny-ftp locally
	go run $(MAIN_PATH)

.PHONY: build
build: ## Build executable binaries for local execution
	go build -ldflags "-s -w" -o $(OUT_PREFIX) $(MAIN_PATH)

.PHONY: build-all
build-all: ## Build executable binaries for all supported OSs and architectures
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o $(OUT_PREFIX).exe $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o $(OUT_PREFIX).macos $(MAIN_PATH)
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o $(OUT_PREFIX).linux-x64 $(MAIN_PATH)
	GOOS=linux GOARCH=arm GOARM=6 go build -ldflags "-s -w" -o $(OUT_PREFIX).linux-arm6 $(MAIN_PATH)
	GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "-s -w" -o $(OUT_PREFIX).linux-arm7 $(MAIN_PATH)
	GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o $(OUT_PREFIX).linux-arm8 $(MAIN_PATH)

.PHONY: count
count: ## Count number of lines of all go codes
	find . -name "*.go" -type f | xargs wc -l

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
