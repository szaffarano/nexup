REPO 		= github.com/szaffarano
PACKAGE  	= nexup
DATE    	?= $(shell date +%FT%T%z)
VERSION 	?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
			cat $(CURDIR)/.version 2> /dev/null || echo v0)

PLATFORMS 	:= linux/amd64 windows/amd64 windows/386
DIST_DIR	:= dist

OS = $(shell echo $@ | cut -d"/" -f1)
ARCH = $(shell echo $@ | cut -d"/" -f2)
EXT = $(shell [ "$(OS)" = "windows" ] && echo -n ".exe")

TRUSTSTORES ?= $(shell cat $(CURDIR)/cacerts 2>/dev/null)

release: $(PLATFORMS)

clean:
	rm -rf $(DIST_DIR)

$(PLATFORMS): clean
	GOOS=$(OS) GOARCH=$(ARCH) go build \
		-tags release \
		-ldflags '-X $(REPO)/$(PACKAGE)/cmd.Version=$(VERSION) -X $(REPO)/$(PACKAGE)/cmd.BuildDate=$(DATE) -X "$(REPO)/$(PACKAGE)/repository.Truststores=$(TRUSTSTORES)"' \
		-o '$(DIST_DIR)/$(PACKAGE)-$(VERSION)-$(OS)-$(ARCH)${EXT}' main.go

.PHONY: release $(PLATFORMS)