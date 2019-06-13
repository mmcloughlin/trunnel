PKG = github.com/mmcloughlin/trunnel
CMD = $(PKG)/cmd/trunnel
GITSHA = `git rev-parse --short HEAD`
LDFLAGS = "-X $(PKG)/meta.GitSHA=$(GITSHA)"

SRC = $(shell find . -type f -name '*.go')
SRC_EXCL_GEN = $(shell find . -type f -name '*.go' -not -name 'gen-*.go')

.PHONY: install
install:
	go install -a -ldflags $(LDFLAGS) $(CMD)

.PHONY: generate
generate: tools
	go generate -x ./...

.PHONY: readme
readme:
	embedmd -w README.md

.PHONY: lint
lint:
	golangci-lint run

.PHONY: imports
imports:
	gofumports -w -local $(PKG) $(SRC)

.PHONY: fmt
fmt:
	gofmt -w -s $(SRC)

.PHONY: cloc
cloc:
	cloc $(SRC_EXCL_GEN)

docs/manual.html: ref/trunnel/doc/trunnel.md
	mkdir -p docs
	markdown  $^ > $@

.PHONY: tools
tools:
	go get -u \
		github.com/mna/pigeon \
		github.com/campoy/embedmd \
		mvdan.cc/gofumpt/gofumports
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(GOPATH)/bin v1.17.1
