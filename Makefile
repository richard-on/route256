
build-all:
	cd checkout && make build
	cd loms && make build
	cd notification && make build

run-all: build-all
	docker compose up -d --force-recreate --build
	cd checkout && exec ./migration.sh
	cd loms && exec ./migration.sh

migrate:
	cd checkout && exec ./migration.sh
	cd loms && exec ./migration.sh

precommit:
	cd checkout && go mod tidy && make precommit
	cd loms && go mod tidy && make precommit
	cd notification && go mod tidy && make precommit
