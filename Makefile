GOARCH = amd64

UNAME = $(shell uname -s)

ifndef OS
	ifeq ($(UNAME), Linux)
		OS = linux
	else ifeq ($(UNAME), Darwin)
		OS = darwin
	endif
endif

.DEFAULT_GOAL := all

all: fmt  build  test start 

build:
	GOOS=$(OS) GOARCH="$(GOARCH)" go build -o vault/plugins/vault-plugin-secrets-datadog cmd/vault-plugin-secrets-datadog/main.go

start:
	vault server -dev -dev-root-token-id=root -dev-plugin-dir=./vault/plugins -log-level=DEBUG -dev-listen-address="127.0.0.1:8200"

test:
	go test -v ./...

enable:
	vault secrets enable -path=datadog vault-plugin-secrets-datadog

clean:
	rm -f ./vault/plugins/vault-plugin-secrets-datadog

fmt:
	go fmt $$(go list ./...)

setup:	enable
	vault write datadog/config  api_key=${DATADOG_API_KEY} app_key=${DATADOG_APP_KEY} api_key_id=${DATADOG_API_KEY_ID} app_key_id=${DATADOG_APP_KEY_ID}
	vault write datadog/roles/test app_key_scopes=incident_read,usage_read max_ttl=3h ttl=2h
	vault read datadog/roles/test

.PHONY: build clean fmt start  enable test setup