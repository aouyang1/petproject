# Go parameters
GOCMD=go
APPNAME=petfinder-updator
TEST_FOUND = $(shell docker ps -a --filter name=$(APPNAME)-test | grep -v CONTAINER | wc -l | xargs)
RUN_FOUND:=$(shell docker ps -a --filter name=$(APPNAME)-run | grep -v CONTAINER | wc -l | xargs)


all: build test
build:
	docker build -t $(APPNAME) .

ifeq ($(TEST_FOUND),1)
test:
	docker stop $(APPNAME)-test
	docker rm $(APPNAME)-test
	docker run --name $(APPNAME)-test $(APPNAME) go test -v ./...
else
test:
	docker run --name $(APPNAME)-test $(APPNAME) go test -v ./...
endif

ifeq ($(RUN_FOUND),1)
run:
	docker stop $(APPNAME)-run
	docker rm $(APPNAME)-run
	docker run --name $(APPNAME)-run -e PETFINDER_API_KEY=$(PETFINDER_API_KEY) -p 8080:8080 $(APPNAME)
else
run:
	docker run --name $(APPNAME)-run -e PETFINDER_API_KEY=$(PETFINDER_API_KEY) -p 8080:8080 $(APPNAME)
endif
