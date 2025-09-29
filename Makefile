# https://maex.me/2018/02/dont-fear-the-makefile/

kitsune:  go-imports 

help: ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/^\([^:]*\):.*##/\1 : /' -e 's/##//'

go-imports:
	go tool goimports -w .

clean:
	go clean -cache -i

nilcheck:
	# go run github.com/uber-go/nilaway@latest ./...

lint:  
	go vet ./...
	go tool staticcheck ./...
	# golangci-lint run ./...

audit:
	go list -json -deps ./... | go tool github.com/sonatype-nexus-community/nancy sleuth --loud

upgrade-deps:
	go get -u ./...
	go mod tidy
	go tool go install gotest.tools/gotestsum@latest ./...

sec: audit
	go tool gosec ./...
	# go run golang.org/x/vuln/cmd/govulncheck@latest ./...

test: go-imports 
	go tool gotestsum ./...
