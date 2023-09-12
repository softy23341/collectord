VERSION ?= $(shell verctl current)
GO := GO15VENDOREXPERIMENT=1 go
GLIDEPATH := $(GOPATH)/src/github.com/Masterminds/glide
GLIDE := $(GLIDEPATH)/glide

build: glide verctl
	$(GLIDE) install
	$(GO) build -o bin/npusherd -ldflags "-X main.version=$(VERSION)" cmd/npusherd/npusherd.go
	$(GO) build -o bin/npusherd-cli -ldflags "-X main.version=$(VERSION)" cmd/npusherd-cli/npusherd-cli.go

glide:
	if [ ! -d $(GLIDEPATH) ]; then git clone https://github.com/Masterminds/glide.git $(GLIDEPATH); fi
	if [ ! -x $(GLIDE) ]; then cd $(GLIDEPATH) && GO15VENDOREXPERIMENT=1 make build; fi

verctl:
	go get gorep.qubeidea.net/verctl

.PHONY: build glide verctl
