# https://maex.me/2018/02/dont-fear-the-makefile/

build: go-imports test ## Build the project and update code coverage in README
	@echo "Building project and updating coverage..."
	@./update_coverage.sh
	@echo "Build complete!"
	
help: ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/^\([^:]*\):.*##/\1 : /' -e 's/##//'

go-imports:
	go tool goimports -w .

clean:
	go clean -cache -i

nilcheck:
	go tool nilaway ./...

lint:  
	go vet ./...
	go tool staticcheck ./...
	golangci-lint run ./...

sec: audit
	go tool gosec  ./... 
	go tool govulncheck  ./...

audit:
	# Error: An error occurred: [401 Unauthorized] error accessing OSS Index
	# go list -json -deps ./... | go tool github.com/sonatype-nexus-community/nancy sleuth --loud

upgrade-deps:
	go get -u ./...
	go mod tidy
	go tool gotestsum  ./...

alltest:  go-imports lint sec nilcheck test

test:
	go tool gotestsum ./...

coverage: ## Generate and display code coverage
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	@rm -f coverage.out
