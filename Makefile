PREFIX=/go
CMD=crc-driver-libvirt
DESCRIBE=$(shell git describe --tags)
CONTAINER_RUNTIME ?= podman

TARGETS=$(addprefix $(CMD)-, centos8 ubuntu20.04)

build: $(TARGETS)

local:
	go build -i -v -o crc-driver-libvirt-local ./cmd/machine-driver-libvirt

$(CMD)-%: Containerfile.%
	${CONTAINER_RUNTIME} rmi -f $@ >/dev/null  2>&1 || true
	${CONTAINER_RUNTIME} rm -f $@-extract > /dev/null 2>&1 || true
	echo "Building binaries for $@"
	${CONTAINER_RUNTIME} build -t $@ -f $< .
	${CONTAINER_RUNTIME} create --name $@-extract $@ sh
	${CONTAINER_RUNTIME} cp $@-extract:$(PREFIX)/bin/$(CMD) ./
	mv ./$(CMD) ./$@
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
lint:
	golangci-lint run

.PHONY: vendorcheck
vendorcheck:
	./verify-vendor.sh

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor
