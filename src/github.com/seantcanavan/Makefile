VERSION_FILE="main/version.no"
BUILD_NUMBER:=$(shell cat ${VERSION_FILE})
INCREMENT_NUMBER=1
NEW_BUILD_NUMBER=$(shell echo $$(( $(BUILD_NUMBER) + $(INCREMENT_NUMBER) )) )
PACKAGES=`go list ./... | grep -v /vendor/`

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

test:
	for p in $(PACKAGES); do \
		go test ../../$$p; \
	done

version-update:
	@echo "Current build number: $(BUILD_NUMBER)"
	@echo "New build number: $(NEW_BUILD_NUMBER)"
	@echo $(NEW_BUILD_NUMBER) > $(VERSION_FILE)
