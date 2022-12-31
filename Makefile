include .env
export $(shell sed 's/=.*//' .env)

.PHONY: ui
dev_ui:
	cd ui \
		&& npm start
dev_server:
	go install github.com/cosmtrek/air@latest
	air

gen:
	go get github.com/99designs/gqlgen
	go run github.com/99designs/gqlgen generate

drun:
	docker build -t ledger .
	docker run -ti \
		-e OIDC_ISSUER_URL \
		-e OIDC_CLIENT_ID \
		-v $(PWD)/data:/app/data \
		-p 8080:8080 ledger

HELM_DOCS = $(shell pwd)/bin/helm-docs
helm-docs: ## Download helm-docs locally if necessary.
	$(call go-get-tool,$(HELM_DOCS),github.com/norwoodj/helm-docs/cmd/helm-docs@v1.6.0)


# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef