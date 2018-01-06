# Go parameters
GOCMD=go
APPNAME=petfinder-updator

all: build test
build:
	docker build -t $(APPNAME) .
test:
	docker stop $(APPNAME)-test
	docker rm $(APPNAME)-test
	docker run --name $(APPNAME)-test $(APPNAME) go test -v ./...
run:
	docker stop $(APPNAME)-run
	docker rm $(APPNAME)-run
	docker run --name $(APPNAME)-run $(APPNAME)
