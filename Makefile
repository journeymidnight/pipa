.PHONY: build
URL = github.com/journeymidnight
REPO = pipa
WORKDIR = /work
BUILDROOT = rpm-build
BUILDDIR = $(WORKDIR)/$(BUILDROOT)/BUILD/$(REPO)
export GO111MODULE=on
export GOPROXY=https://goproxy.cn

build:
    docker run --rm -v $(PWD):$(BUILDDIR) -w $(BUILDDIR) journeymidnight/pipa bash -c 'make build_internal'

build_internal:
    go build $(URL)/$(REPO)

pkg:
    go build $(URL)/$(REPO) && cd package && bash rpmbuild.sh $(REPO)

env:

