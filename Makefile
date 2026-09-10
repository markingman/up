# Makefile for local development

-include ./.env

.DEFAULT_GOAL := help
.PHONY: help
NAME=up

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ": ## "}; {printf "\033[36m%-28s\033[0m %s\n", $$1, $$2}' | sed 's/Makefile://g'

build: ## Build a Docker image for local development
	@docker build -t $(NAME) .

run: ## Run the Docker image
	@docker run -v `pwd`/disk:/app/data --name $(NAME) --env SMTP_PASSWD="${SMTP_PASSWD}" $(NAME)

start: ## Start Docker container to run tests (if container built and stopped)
	@docker container start $(NAME)

stop: ## Stop current container (if running)
	@docker stop $(NAME)

ssh: ## SSH to Docker container
	@docker exec -it $(NAME) sh

clean: ## Clean up
	@docker stop $(NAME)
	@docker rm $(NAME)
