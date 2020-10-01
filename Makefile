
# install with go get -u github.com/go-bindata/go-bindata/...
.PHONY: gen
gen:
	go-bindata -prefix assets -pkg assets -o internal/generated/assets/assets.go assets

.PHONY: build
build:
	go build ./cmd/md3pdf

.PHONY: clean
clean:
	rm -rf internal/generated
	rm -f md3pdf