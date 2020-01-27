
BUILD	= go build
CLEAN	= go clean -i -cache
TEST	= go test
INSTALL = go install
ADDPKG  = dep ensure
BINARY_NAME = logour

all: test build
build: clean deps test
	$(BUILD) -o $(BINARY_NAME) -v
test:
	$(TEST) -v ./...
clean:
	$(CLEAN)
	rm -f $(BINARY_NAME)
run: build
	./$(BINARY_NAME)
deps:
	$(ADDPKG)
install: build
	$(INSTALL)
