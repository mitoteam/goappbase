# ------ USAGE EXAMPLE -------
# EXECUTABLE_NAME := exec_name
# SUBMODULES := ./internal/dhtml
#
# include internal/goappbase/Makefile.inc.mk
# ----------------------------


# SET DEFAULT VALUES
ifeq ($(EXECUTABLE_NAME),)
EXECUTABLE_NAME := executable_name_not_set
endif

# directory to create distribution archives
ifeq ($(DIST_DIR),)
DIST_DIR := dist
endif

# internal submodules base path
ifeq (${INTERNAL_SUBMODULES_PATH},)
INTERNAL_SUBMODULES_PATH := ./internal
endif

# submodules to run tests for
SUBMODULES := ${SUBMODULES} ${INTERNAL_SUBMODULES_PATH}/mttools ${INTERNAL_SUBMODULES_PATH}/goappbase

#set default target
.DEFAULT_GOAL := build-all


ifeq (${OS},Windows_NT)
	BUILD_TIME := $(shell powershell "Get-Date -Format 'yyyy-MM-dd HH:mm:ss'")
else
	BUILD_TIME := $(shell date +"%Y-%m-%d %H:%M:%S")
endif

APP_VERSION := $(file < VERSION)
APP_COMMIT := $(shell git rev-list -1 HEAD)
LD_FLAGS := "-X 'github.com/mitoteam/goappbase.BuildVersion=${APP_VERSION}' -X 'github.com/mitoteam/goappbase.BuildCommit=${APP_COMMIT}' -X 'github.com/mitoteam/goappbase.BuildTime=${BUILD_TIME}'"

fn_GO_BUILD = GOOS=$(1) GOARCH=$(2) go build -o ${DIST_DIR}/$(3) -ldflags=${LD_FLAGS} main.go ;\
7z a ${DIST_DIR}/${EXECUTABLE_NAME}-${APP_VERSION}-$(4).7z -mx9 ./${DIST_DIR}/$(3)


.PHONY: build-all
build-all:: build-linux64 build-linux32 build-windows32 build-windows64
	rm -f ${DIST_DIR}/${EXECUTABLE_NAME}
	rm -f ${DIST_DIR}/${EXECUTABLE_NAME}.exe


.PHONY: build-windows32
build-windows32: clean tests ${DIST_DIR}
	$(call fn_GO_BUILD,windows,386,${EXECUTABLE_NAME}.exe,win64)

.PHONY: build-windows64
build-windows64: clean tests ${DIST_DIR}
	$(call fn_GO_BUILD,windows,amd64,${EXECUTABLE_NAME}.exe,win64)

.PHONY: build-linux32
build-linux32: clean tests ${DIST_DIR}
	$(call fn_GO_BUILD,linux,386,${EXECUTABLE_NAME},linux32)

.PHONY: build-linux64
build-linux64: clean tests ${DIST_DIR}
	$(call fn_GO_BUILD,linux,amd64,${EXECUTABLE_NAME},linux64)


# Run all tests in root module and in known submodules
.PHONY: tests
tests::
	clear
	go test ./... $(SUBMODULES)


${DIST_DIR}:
	mkdir ${DIST_DIR}


.PHONY: version
version:
	@echo Version from 'VERSION' file: ${APP_VERSION}


.PHONY: clean
clean:
	rm -rf ${DIST_DIR}
