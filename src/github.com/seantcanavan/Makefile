VERSION_FILE="main/version.no"
BUILD_NUMBER:=$(shell cat ${VERSION_FILE})
INCREMENT_NUMBER=1
NEW_BUILD_NUMBER=$(shell echo $$(( $(BUILD_NUMBER) + $(INCREMENT_NUMBER) )) )
PACKAGES=config main updater



updater:
	echo "building the updater"

client:
	echo "building the client"

version-update:
	@echo "Current build number: $(BUILD_NUMBER)"
	@echo "New build number: $(NEW_BUILD_NUMBER)"
	@echo $(NEW_BUILD_NUMBER) > $(VERSION_FILE)

format:
	gofmt -w config
	gofmt -w main
	gofmt -w updater
