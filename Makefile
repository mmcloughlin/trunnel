SRC = $(shell find . -type f -name '*.go')
SRC_EXCL_GEN = $(shell find . -type f -name '*.go' -not -name 'gen-*.go')

generate: deps
	go generate -x ./...

lint:
	gometalinter --config=.gometalinter.json ./...

imports:
	goimports -w -local github.com/mmcloughlin/trunnel $(SRC)

fmt:
	gofmt -w -s $(SRC)

cloc:
	cloc $(SRC_EXCL_GEN)

deps:
	go get -u github.com/mna/pigeon
