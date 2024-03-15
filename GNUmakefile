default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
	go test ./... -sweep empty

# Build locally
.PHONY: local-build
local-build:
	go install .
