# Go parameters
GOCMD=go
APPNAME=petfinder-updator
TEST_FOUND = $(shell docker ps -a --filter name=$(APPNAME)-test | grep -v CONTAINER | wc -l | xargs)

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

run:
	docker-compose up

stop:
	docker-compose down
	
