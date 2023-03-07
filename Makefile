
build-all:
	cd checkout && make build
	cd loms && make build
	cd notification && make build

run-all: build-all
	docker compose up --force-recreate --build

precommit:
	cd checkout && go mod tidy && make precommit
	cd loms && go mod tidy && make precommit
	cd notification && go mod tidy && make precommit
