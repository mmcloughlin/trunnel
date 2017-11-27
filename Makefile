SRC = $(shell find . -type f -name '*.go')

imports:
	goimports -w -local github.com/mmcloughlin/trunnel $(SRC)
