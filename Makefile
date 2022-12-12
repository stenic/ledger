include .env
export $(shell sed 's/=.*//' .env)

.PHONY: ui
ui:
	cd ui \
		&& npm start

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
