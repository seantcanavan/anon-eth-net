VERSION_FILE="main/version.no"
BUILD_NUMBER:=$(shell cat ${VERSION_FILE})
INCREMENT_NUMBER=1
NEW_BUILD_NUMBER=$(shell echo $$(( $(BUILD_NUMBER) + $(INCREMENT_NUMBER) )) )
PACKAGES=`go list ./... | grep -v /vendor/`

default: all

all: clean format version-update install;

check: vet test

install:
	for p in $(PACKAGES); do \
		go install ../../$$p; \
	done

format:
	for p in $(PACKAGES); do \
		gofmt -w ../../$$p; \
	done

clean:
	for p in $(PACKAGES); do \
		go clean ../../$$p; \
	done
# cleanup temporary files created after test
	find . -name "*.log" -type f -delete
	find . -name "*.rep" -type f -delete
	find . -name "*.tar" -type f -delete

test:
	for p in $(PACKAGES); do \
		go test ../../$$p; \
	done

vet:
	for p in $(PACKAGES); do \
		go vet ../../$$p; \
	done

version-update:
	@echo "Current build number: $(BUILD_NUMBER)"
	@echo "New build number: $(NEW_BUILD_NUMBER)"
	@echo $(NEW_BUILD_NUMBER) > $(VERSION_FILE)
