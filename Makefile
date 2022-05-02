.PHONY: all build

help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)


dev: ## Start the containers and start outputting logs
	docker-compose up

build: ## Force a full rebuild of the development containers
	docker-compose build --no-cache

clean: ## Remove any build artifacts
	docker container rm -f go-proxy-yourself
	docker image rm -f go-proxy-yourself_go-proxy-yourself
	rm go-proxy-yourself

