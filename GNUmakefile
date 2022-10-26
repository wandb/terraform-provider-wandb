default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

build-and-install:
	go build
	mv ./terraform-provider-wandb ~/.terraform.d/plugins/registry.terraform.io/local/wandb/0.0.1/darwin_arm64