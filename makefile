
BUILD	= go build
CLEAN	= go clean -i -cache
TEST	= go test
ADDPKG	= go get
INSTALL = go install
BINARY_NAME = logour

all: test build
build:
	$(BUILD) -o $(BINARY_NAME) -v
test:
	$(TEST) -v ./...
clean:
	$(CLEAN)
	rm -f $(BINARY_NAME)
run: build
	./$(BINARY_NAME)
deps:
	$(ADDPKG) -v
install: clean test
	$(INSTALL)
