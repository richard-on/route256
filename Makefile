
build-all:
	cd checkout && make build
	cd loms && make build
	cd notification && make build

run-all: build-all
	docker compose up --force-recreate --build

precommit:
	cd checkout && make precommit
	cd loms && make precommit
	cd notification && make precommit
