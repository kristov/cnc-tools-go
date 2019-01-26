GOPATH := $(shell pwd)
GO := /opt/go/bin/go
SRCDIR := src
BINDIR := bin

all: bin/cnc-stl-view

$(BINDIR)/%: $(SRCDIR)/%.go
	@GOPATH=$(GOPATH) $(GO) build -o $@ $<

deps:
	@GOPATH=$(GOPATH) $(GO) get -v github.com/hschendel/stl
	@GOPATH=$(GOPATH) $(GO) get -v github.com/go-gl/gl/v2.1/gl
	@GOPATH=$(GOPATH) $(GO) get -v github.com/go-gl/glfw/v3.2/glfw
	@GOPATH=$(GOPATH) $(GO) get -v github.com/go-gl/mathgl/mgl32
	#@GOPATH=$(GOPATH) $(GO) get github.com/deadsy/sdfx/sdf

clean:
	rm -f bin/*
