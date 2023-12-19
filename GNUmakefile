TEST?=$$(go list ./... |grep -v 'vendor')
default: build

# Run acceptance tests
.PHONY: testacc

test: lint
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=8

testacc: lint
	TF_ACC=1 go test -v -count=1 -timeout 600m $(TEST_FILEPATH)

generate:
	find examples -name "*.tf" -exec terraform fmt {} \;
	go generate ./...

lint:
	golangci-lint run

lintWithFix:
	golangci-lint run --fix

build: lint
	install 

install:
	go install .

submodules:
	@git submodule sync
	@git submodule update --init --recursive
	@git config core.hooksPath githooks
	@git config submodule.recurse true
