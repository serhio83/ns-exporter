# основные пути
GOPATH			:=	${shell pwd}
BINPATH			=	$(GOPATH)/bin

# основные команды
GOCMD			=	go

# параметры команд
GODEP			=	$(GOCMD) get
GOBUILD			=	$(GOCMD) build

GOTEST			=	$(GOCMD) test
GOINSTALL		=	$(GOCMD) install

export GOPATH

MAIN_PKGS 		:=	github.com/serhio83/ns-exporter

LIBS_PKGS		:=

DEPS_PKGS 		:=	github.com/golang/protobuf/proto \
                    github.com/golang/protobuf/protoc-gen-go \
                    github.com/docker/docker/client \
                    github.com/docker/docker/api/types \
                    github.com/prometheus/client_golang/prometheus

TEST_PKGS		:=	$(LIBS_PKGS) $(MAIN_PKGS)

DEPS_LIST		=	$(foreach int, $(DEPS_PKGS), $(int)_deps)
BUILD_LIST		=	$(foreach int, $(MAIN_PKGS), $(int)_build)

TEST_LIST		=	$(foreach int, $(TEST_PKGS), $(int)_test)
INSTALL_LIST	=	$(foreach int, $(MAIN_PKGS), $(int)_install)

.PHONY:			$(DEPS_LIST) $(TEST_LIST) $(BUILD_LIST) $(INSTALL_LIST)

all:			build
deps:			$(DEPS_LIST)
test:			$(TEST_LIST)
build:			$(BUILD_LIST)
install:		$(INSTALL_LIST)

$(DEPS_LIST): %_deps:
	$(GODEP) $*
$(BUILD_LIST): %_build:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINPATH)/$(shell basename $*) $*

$(TEST_LIST): %_test:
	$(GOTEST) $*
$(INSTALL_LIST): %_install:
	$(GOINSTALL) $*
