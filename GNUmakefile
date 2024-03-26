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

# Generate docs
.PHONY: docs
docs:
	@which tfplugindocs &>/dev/null || (echo "install tfplugindocs (https://github.com/hashicorp/terraform-plugin-docs)"; exit 1)
	tfplugindocs generate .
