PKG = github.com/mmcloughlin/trunnel
CMD = $(PKG)/cmd/trunnel
GITSHA = `git rev-parse --short HEAD`
LDFLAGS = "-X $(PKG)/meta.GitSHA=$(GITSHA)"

SRC = $(shell find . -type f -name '*.go')
SRC_EXCL_GEN = $(shell find . -type f -name '*.go' -not -name 'gen-*.go')

.PHONY: install
install:
	go install -ldflags $(LDFLAGS) $(CMD)

.PHONY: generate
generate: deps
	go generate -x ./...

.PHONY: lint
lint:
	gometalinter --config=.gometalinter.json ./...

.PHONY: imports
imports:
	goimports -w -local $(PKG) $(SRC)

.PHONY: fmt
fmt:
	gofmt -w -s $(SRC)

.PHONY: cloc
cloc:
	cloc $(SRC_EXCL_GEN)

.PHONY: deps
deps:
	go get -u github.com/mna/pigeon
