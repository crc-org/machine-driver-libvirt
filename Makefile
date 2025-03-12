CMD=crc-driver-libvirt
DESCRIBE=$(shell git describe --tags)
VERSION=$(shell grep DriverVersion pkg/libvirt/constants.go | awk '{ print $$3 }' | tr -d \")
CONTAINER_RUNTIME ?= podman

TARGETS=$(addprefix $(CMD)-, stream9)

build: $(TARGETS)

TOOLS_DIR := tools
include tools/tools.mk

local:
	go build -v -o crc-driver-libvirt-local ./cmd/machine-driver-libvirt

$(CMD)-%: Containerfile.%
	${CONTAINER_RUNTIME} rmi -f $@ >/dev/null  2>&1 || true
	${CONTAINER_RUNTIME} rm -f $@-extract > /dev/null 2>&1 || true
	echo "Building binaries for $@"
	${CONTAINER_RUNTIME} build -t $@ -f $< .
	${CONTAINER_RUNTIME} create --name $@-extract $@ sh
	${CONTAINER_RUNTIME} cp $@-extract:/go/bin/$(CMD) ./$@
	${CONTAINER_RUNTIME} rm $@-extract || true
	${CONTAINER_RUNTIME} rmi $@ || true

clean:
	rm -f ./$(CMD)-*

release: build
	@echo "Paste the following into the release page on github and upload the binaries..."
	@echo ""
	@for bin in $(CMD)-* ; do \
	    target=$$(echo $${bin} | cut -f5- -d-) ; \
	    md5=$$(md5sum $${bin}) ; \
	    echo "* $${target} - md5: $${md5}" ; \
	    echo '```' ; \
	    echo "  curl -L https://github.com/dhiltgen/docker-machine-kvm/releases/download/$(DESCRIBE)/$${bin} > /usr/local/bin/$(CMD) \\ " ; \
	    echo "  chmod +x /usr/local/bin/$(CMD)" ; \
	    echo '```' ; \
	done

.PHONY: validate
validate: test lint vendorcheck

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint: $(TOOLS_BINDIR)/golangci-lint
	${TOOLS_BINDIR}/golangci-lint run

.PHONY: vendorcheck
vendorcheck:
	./verify-vendor.sh

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

.PHONY: spec test-rpmbuild
spec: crc-driver-libvirt.spec

test-rpmbuild: spec
	${CONTAINER_RUNTIME} build -f Containerfile.rpmbuild .

%.spec: %.spec.in $(TOOLS_BINDIR)/gomod2rpmdeps
	@"$(TOOLS_BINDIR)"/gomod2rpmdeps | sed -e '/__BUNDLED_REQUIRES__/r /dev/stdin' \
					   -e '/__BUNDLED_REQUIRES__/d' \
					   -e 's/__VERSION__/'$(VERSION)'/g' \
				       $< >$@
