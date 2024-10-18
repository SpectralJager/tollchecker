up:
	@clear
	@docker compose build
	@docker compose up

clear:
	# @docker rmi $(docker images -f "dangling=true" -q)

.PHONY: up