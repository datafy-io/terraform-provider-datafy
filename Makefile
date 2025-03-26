default: fmt install

.PHONY: build
build:
	go build -v ./...

.PHONY: install
install: build
	go install -v ./...

.PHONY: fmt
fmt:
	gofmt -s -w -e .

.PHONY: test
test:
	go test -v -cover -timeout=120s -parallel=10 ./...


