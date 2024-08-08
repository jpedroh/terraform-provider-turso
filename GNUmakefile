default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

update_sdk:
	go run github.com/ogen-go/ogen/cmd/ogen --target internal/client -package client --clean internal/client/openapi.json