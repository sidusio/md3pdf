# This file requires the $GOPATH variable to be set in the environment

# Finds the GOBIN directory
gobin := $(shell (gobin="$(shell go env GOBIN)"; [[ -n "$$gobin" ]] && echo $$gobin || echo "$(shell go env GOPATH)/bin"))

$(gobin)/go-bindata:
	@echo "You need to download 'go-bindata' to generate the assets"
	@read -p "Is it ok to run 'go get -u github.com/go-bindata/go-bindata/...'? (y/N) " ANSWER; \
	if [[ $$ANSWER = "Y" || $$ANSWER = "y" ]]; then \
			echo "Installing 'go-bindata'..."; \
			go get -u github.com/go-bindata/go-bindata/...; \
		else \
			echo "go-bindata will not be installed"; exit 1; \
	fi

.PHONY: gen
gen: $(gobin)/go-bindata
	go-bindata -prefix assets -pkg assets -o internal/generated/assets/assets.go assets

.PHONY: build
build:
	go build ./cmd/md3pdf

.PHONY: clean
clean:
	rm -rf internal/generated
	rm -f md3pdf

.PHONY: all
all: clean gen build