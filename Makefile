GOPATH_DIR=`go env GOPATH`

.PHONY: test
test:
	go test -count 2 -v -race ./...

.PHONY: lint
lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH_DIR)/bin v1.50.1
	$(GOPATH_DIR)/bin/golangci-lint run ./...
	$(GOPATH_DIR)/bin/golangci-lint run --build-tags example ./_examples/basic
	$(GOPATH_DIR)/bin/golangci-lint run --build-tags example ./_examples/configfile
	$(GOPATH_DIR)/bin/golangci-lint run --build-tags example ./_examples/gamepad_in_browser
	$(GOPATH_DIR)/bin/golangci-lint run --build-tags example ./_examples/modkeys
	$(GOPATH_DIR)/bin/golangci-lint run --build-tags example ./_examples/scroll
	$(GOPATH_DIR)/bin/golangci-lint run --build-tags example ./_examples/simulateinput
	$(GOPATH_DIR)/bin/golangci-lint run --build-tags example ./_examples/action_released
	@echo "everything is OK"
