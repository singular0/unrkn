GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

CMD_NAME=unrkn
    
all: build

build:
	$(GOBUILD) ./cmd/$(CMD_NAME)

clean:
	$(GOCLEAN)
	rm -f $(CMD_NAME)
