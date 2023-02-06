default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

generate:
	go generate ./...

install:
	go install .

test:
	go test -count=1 -parallel=4 ./...

submodules:
	@git submodule sync
	@git submodule update --init --recursive
	@git config core.hooksPath githooks