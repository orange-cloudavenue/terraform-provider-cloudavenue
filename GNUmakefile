default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

generate:
	go install github.com/FrangipaneTeam/tf-doc-extractor@latest
	go generate -run "tf-doc-extractor" ./...
	# golang 1.20 feature
	go generate -skip "tf-doc-extractor" ./...

install:
	go install .

test:
	go test -count=1 -parallel=4 ./...

submodules:
	@git submodule sync
	@git submodule update --init --recursive
	@git config core.hooksPath githooks
	@git config submodule.recurse true
