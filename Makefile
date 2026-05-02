.PHONY: build test install cross-build clean validate-contracts generate-contracts check-contracts

BINARY_NAME = workflow-plugin-okta
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -ldflags "-X main.version=$(VERSION)"
INSTALL_DIR ?= data/plugins/$(BINARY_NAME)
PLATFORMS = linux/amd64 linux/arm64 darwin/amd64 darwin/arm64
LOCAL_BIN = $(CURDIR)/bin
WFCTL_VERSION ?= v0.20.5
WFCTL = $(LOCAL_BIN)/wfctl

build:
	GOPRIVATE=github.com/GoCodeAlone/* go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/$(BINARY_NAME)

test:
	GOPRIVATE=github.com/GoCodeAlone/* go test ./... -v -race

install: build
	mkdir -p $(DESTDIR)/$(INSTALL_DIR)
	cp bin/$(BINARY_NAME) $(DESTDIR)/$(INSTALL_DIR)/
	cp plugin.json $(DESTDIR)/$(INSTALL_DIR)/
	cp plugin.contracts.json $(DESTDIR)/$(INSTALL_DIR)/

cross-build:
	@mkdir -p bin
	@for platform in $(PLATFORMS); do \
		os=$${platform%%/*}; \
		arch=$${platform##*/}; \
		output=bin/$(BINARY_NAME)-$${os}-$${arch}; \
		echo "Building $${output}..."; \
		CGO_ENABLED=0 GOOS=$${os} GOARCH=$${arch} GOPRIVATE=github.com/GoCodeAlone/* \
			go build $(LDFLAGS) -o $${output} ./cmd/$(BINARY_NAME); \
	done

$(WFCTL):
	mkdir -p $(LOCAL_BIN)
	GOBIN=$(LOCAL_BIN) go install github.com/GoCodeAlone/workflow/cmd/wfctl@$(WFCTL_VERSION)

# plugin.contracts.json is generated from plugin.json capabilities so strict
# contract coverage stays in sync with the advertised module and step types.
generate-contracts:
	jq '{version:"1",contracts:(([.capabilities.moduleTypes[]|{kind:"module",type:.,mode:"strict",config:"workflow.plugins.okta.v1.ProviderConfig"}]+[.capabilities.stepTypes[]|{kind:"step",type:.,mode:"strict",config:"workflow.plugins.okta.v1.OktaStepConfig",input:"workflow.plugins.okta.v1.OktaStepInput",output:"workflow.plugins.okta.v1.OktaStepOutput"}]))}' plugin.json > plugin.contracts.json

check-contracts: generate-contracts
	git diff --exit-code -- plugin.contracts.json

validate-contracts: $(WFCTL)
	$(WFCTL) plugin validate --file plugin.json --strict-contracts

clean:
	rm -rf bin/
