build:
	go build -ldflags="-X kerosenelabs.com/espresso/core/service.CommitSha=$(shell git rev-parse HEAD) -X kerosenelabs.com/espresso/core/service.Version=$(VERSION)"