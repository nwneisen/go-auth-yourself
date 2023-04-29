.PHONY: all build dev clean vendor

COMPOSE=docker compose -f build/package/docker-compose.yml

help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

dev: certs ## Start the containers and start outputting logs
	scripts/setup-dev-config
	$(COMPOSE) up

build: ## Force a full rebuild of the development containers
	$(COMPOSE) build --no-cache

clean: ## Remove any build artifacts
	$(COMPOSE) down
	$(COMPOSE) rm -f -s -v
	docker image rm -f package-go-proxy-yourself 2>/dev/null
	docker image rm -f httpd 2>/dev/null
	rm -rf bin

mod: ## Update dependencies
	go get -u ./...
	go mod tidy
	go mod vendor

certs: ## Generate new certs
	scripts/generate-certs
