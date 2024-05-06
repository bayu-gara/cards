pkgs = $(shell go list ./... | grep -v vendor | grep -v mock)

# go build command
gobuild:
	@go build -v -o cards cmd/cards/*.go

# go run command
# gorun default to rest
gorun:
	make gobuild
	./cards --mode=rest

gotest:
	@echo "RUN TESTING..."
	@go test -v -cover -gcflags=-l -race $(pkgs) -coverprofile coverage.out	