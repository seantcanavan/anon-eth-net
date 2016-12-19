VERSION_FILE="assets/version.no"
BUILD_NUMBER:=$(shell cat ${VERSION_FILE})
INCREMENT_NUMBER=1
NEW_BUILD_NUMBER=$(shell echo $$(( $(BUILD_NUMBER) + $(INCREMENT_NUMBER) )) )
PACKAGES=`go list ./... | grep -v /vendor/`

default: all

clean:
	for p in $(PACKAGES); do \
		go clean ../../$$p; \
	done
# cleanup temporary files created after test
	find . -name "*.log" -type f -delete
	find . -name "*.rep" -type f -delete
	find . -name "*.tar" -type f -delete
	find . -name "*.txt" -type f -delete
	find . -name "*.run" -type f -delete

deps:
	glide install

format:
	for p in $(PACKAGES); do \
		gofmt -w ../../$$p; \
	done

install:
	for p in $(PACKAGES); do \
		go install ../../$$p; \
	done

release:
	for p in $(PACKAGES); do \
		go install ../../$$p; \
	done
	cp -f ../../../bin/main ./bin/anon-eth-net

test:
	for p in $(PACKAGES); do \
		go test -v ../../$$p; \
	done

vet:
	for p in $(PACKAGES); do \
		go vet ../../$$p; \
	done

ubuntu-zip:
	zip -r "anon-eth-net-v0.1.0-ubuntu-x64.zip" assets/ bin/

version-update:
	@echo "Current build number: $(BUILD_NUMBER)"
	@echo "New build number: $(NEW_BUILD_NUMBER)"
	@echo $(NEW_BUILD_NUMBER) > $(VERSION_FILE)

#Aggregate commands

all: clean version-update format vet install test;

prepare-commit: clean version-update format vet install test clean;

test-clean: test clean
