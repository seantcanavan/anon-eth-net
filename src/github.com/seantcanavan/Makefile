VERSION_FILE="main/version.no"
BUILD_NUMBER:=$(shell cat ${VERSION_FILE})
INCREMENT_NUMBER=1
NEW_BUILD_NUMBER=$(shell echo $$(( $(BUILD_NUMBER) + $(INCREMENT_NUMBER) )) )

all:
	go install config
	go install loader
	go install logger
	go install main
	go install profiler
	go install reporter
	go install updater
	go install utils

format:
	gofmt -w config
	gofmt -w loader
	gofmt -w logger
	gofmt -w main
	gofmt -w profiler
	gofmt -w reporter
	gofmt -w updater
	gofmt -w utils

test:
	go test config
	go test loader
	go test logger
	go test main
	go test profiler
	go test reporter
	go test updater
	go test utils

updater:
	echo "building the updater"

client:
	echo "building the client"

version-update:
	@echo "Current build number: $(BUILD_NUMBER)"
	@echo "New build number: $(NEW_BUILD_NUMBER)"
	@echo $(NEW_BUILD_NUMBER) > $(VERSION_FILE)
