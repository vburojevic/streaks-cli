BIN=bin/st

.PHONY: build test lint release-dry integration smoke

build:
	mkdir -p bin
	go build -o $(BIN) ./cmd/streaks-cli

test:
	go test ./...

lint:
	golangci-lint run ./...

release-dry:
	goreleaser release --snapshot --clean

integration:
	STREAKS_CLI_INTEGRATION=1 go test ./internal/cli -run Integration

smoke: build
	@echo "Running smoke checks (requires Streaks installed)..."
	$(BIN) --help >/dev/null
	$(BIN) --agent discover >/dev/null
	$(BIN) --output plain actions list >/dev/null
	$(BIN) --output plain wrappers list >/dev/null
