SHELL := /bin/bash
GOPATH ?= $(shell go env GOPATH)
export PATH := $(GOPATH)/bin:$(PATH)

.PHONY: dev run run-uat run-prd migration mock tidy swagger infra-up infra-down

ENV_FILE ?= .env

run: swagger
	@set -a && source $(ENV_FILE) && set +a && go run ./cmd/backend/main.go

dev: swagger
	@set -a && source $(ENV_FILE) && set +a && air


migration:
	@set -a && source $(ENV_FILE) && set +a && go run ./cmd/migration

mock:
	@set -a && source $(ENV_FILE) && set +a && go run ./cmd/mock

swagger:
	go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/backend/main.go -o docs --parseInternal --parseDependency -q

tidy:
	go mod tidy

infra-up:
	docker compose -f local-service.yaml up -d

infra-down:
	docker compose -f local-service.yaml down


COV_FILE := /tmp/eduid-backend-test-cov.tmp

test:
	@echo -e "\033[1;36m╔══════════════════════════════════════════════════════════════════════════════╗\033[0m"
	@echo -e "\033[1;36m║                          🧪  RUNNING SERVICE TESTS                           ║\033[0m"
	@echo -e "\033[1;36m╚══════════════════════════════════════════════════════════════════════════════╝\033[0m"
	@echo ""
	@go test -coverpkg=./... ./module/... -coverprofile=$(COV_FILE) 2>&1 | grep -E "(PASS|FAIL|ok|FAIL|coverage)" || true
	@echo ""
	@echo -e "\033[1;36m╔══════════════════════════════════════════════════════════════════════════════╗\033[0m"
	@echo -e "\033[1;36m║                       📊  COVERAGE REPORT - PUBLIC FUNCTIONS                 ║\033[0m"
	@echo -e "\033[1;36m╚══════════════════════════════════════════════════════════════════════════════╝\033[0m"
	@echo ""
	@report=$$(go tool cover -func=$(COV_FILE)); \
	echo "$$report" | awk -F'\t+' '/^[^:]*:[0-9]+:[ \t]*[A-Z]/ { \
		val=$$NF; gsub("%","",val); num=val+0; \
		color="\033[0;31m"; icon="❌"; \
		if (num == 100) { color="\033[0;32m"; icon="✅"; } \
		else if (num >= 80) { color="\033[0;33m"; icon="⚠️ "; } \
		sub(/^.*\//, "", $$1); \
		combined = $$1 " " $$2; \
		printf "%s %-70s %s%6s\033[0m\n", icon, combined, color, $$NF \
	}'; \
	echo ""; \
	echo -e "\033[1;36m╔══════════════════════════════════════════════════════════════════════════════╗\033[0m"; \
	echo -e "\033[1;36m║                           📈  TOTAL COVERAGE                                 ║\033[0m"; \
	echo -e "\033[1;36m╚══════════════════════════════════════════════════════════════════════════════╝\033[0m"; \
	echo ""; \
	total_line=$$(echo "$$report" | tail -1); \
	total_val=$$(echo "$$total_line" | awk '{print $$NF}'); \
	total_num=$$(echo "$$total_val" | sed 's/%//'); \
	IFS='|' read -r color icon status <<< "$$(awk -v n="$$total_num" 'BEGIN{ \
		if (n>=80) printf "\033[1;32m|🎉|EXCELLENT"; \
		else if (n>=60) printf "\033[1;33m|👍|GOOD"; \
		else printf "\033[1;31m|⚠️ |NEEDS IMPROVEMENT" \
	}')"; \
	echo -e "$$icon  Status: $$color$$status\033[0m"; \
	echo -e "📊  Total Coverage: $$color$$total_val\033[0m"; \
	echo ""; \
	func_count=$$(echo "$$report" | grep -c "^[^:]*:[0-9]*:"); \
	echo -e "📝  Functions Tested: \033[1;36m$$func_count\033[0m"; \
	echo ""
	@echo -e "\033[1;36m╔══════════════════════════════════════════════════════════════════════════════╗\033[0m"
	@echo -e "\033[1;36m║                              ✨  TEST COMPLETE                               ║\033[0m"
	@echo -e "\033[1;36m╚══════════════════════════════════════════════════════════════════════════════╝\033[0m"
	@rm -f $(COV_FILE)
