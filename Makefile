GOPATH_DIR=`go env GOPATH`

.PHONY: test
test:
	go test -count 2 -v -race ./...

.PHONY: lint
lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH_DIR)/bin v1.50.1
	$(GOPATH_DIR)/bin/golangci-lint run ./...
	@echo "everything is OK"
