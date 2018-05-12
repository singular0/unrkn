
all: test build

build:
	go build ./cmd/unrkn

test:
	go test ./internal/subnet

clean:
	go clean
	rm -f unrkn
