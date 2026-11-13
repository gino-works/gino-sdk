.DEFAULT_GOAL := help

BLUE   := $(shell tput -Txterm setaf 4)
YELLOW := $(shell tput -Txterm setaf 3)
GREEN  := $(shell tput -Txterm setaf 2)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: help

help:
	@echo ""
	@echo "${YELLOW}Usage:${RESET} make ${GREEN}<target>${RESET}"
	@echo ""
	@awk ' \
		/^##@@/ { gsub(/^##@@ */, ""); printf "\n${BLUE}%s${RESET}\n", toupper($$0); next } \
		/^##@/  { gsub(/^##@ */, ""); printf "  ${YELLOW}%s:${RESET}\n", $$0; next } \
		/^## /   { gsub(/^## */, ""); msg = $$0; next } \
		/^[a-zA-Z0-9_-]+:/ { \
			target = $$1; gsub(/:.*/, "", target); \
			if (match($$0, /## */)) { \
				msg = substr($$0, RSTART + RLENGTH); \
			} \
			if (msg != "") { \
				printf "    ${GREEN}%-18s${RESET} %s\n", target, msg; \
				msg = ""; \
			} \
		} \
	' $(MAKEFILE_LIST)
	@echo ""

include Makefile.generation.mk

##@@ Tooling

##@ Format/Lint

.SILENT:
## Run formatters
fmt:
	go tool gofumpt -w .
	go tool gci write .
	go tool golines --base-formatter="gofmt" -w .
	go tool buf format -w

.SILENT:
## Run linters
lint: fmt
	@failed=0; \
	echo "\n\t***lint go files***\n"; \
	golangci-lint run || failed=1; \
	go tool gohawk -enable-all ./... || failed=1; \
	\
	echo "\n\t***lint proto files***\n"; \
	go tool buf lint
	exit $$failed
