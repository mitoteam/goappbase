# EXAMPLE:
# EXECUTABLE_NAME := exec_name
# include goappbase/Makefile.inc.mk


# Expected variables
ifeq ($(EXECUTABLE_NAME),)
EXECUTABLE_NAME := executable_name_not_set #default value if not set
endif

SUBMODULES := $(SUBMODULES) ./mttools ./goappbase

#set default target
.DEFAULT_GOAL := build-all

DIST_DIR := dist

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

HZ := $(shell date)

.PHONY: build-all
build-all: build-linux64 build-linux32 build-windows64
	rm -f ${DIST_DIR}/${EXECUTABLE_NAME}
	rm -f ${DIST_DIR}/${EXECUTABLE_NAME}.exe


.PHONY: build-windows64
build-windows64: clean tests ${DIST_DIR}
	$(call fn_GO_BUILD,windows,amd64,${EXECUTABLE_NAME}.exe,win64)

.PHONY: build-linux32
build-linux32: clean tests ${DIST_DIR}
	$(call fn_GO_BUILD,linux,386,${EXECUTABLE_NAME},linux32)

.PHONY: build-linux64
build-linux64: clean tests ${DIST_DIR}
	$(call fn_GO_BUILD,linux,amd64,${EXECUTABLE_NAME},linux64)

.PHONY: tests
tests::
# all tests in root module and in known submodules
	go test ./... $(SUBMODULES)


${DIST_DIR}:
	mkdir ${DIST_DIR}


.PHONY: version
version:
	@echo Version from 'VERSION' file: ${APP_VERSION}


.PHONY: clean
clean:
	rm -rf ${DIST_DIR}
