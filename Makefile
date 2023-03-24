
test-coverage:
	go test -v ./checkout/... ./loms/... -coverprofile=coverage.out
	go tool cover -func=./coverage.out | grep "total";
	go tool cover -html=coverage.out

build-all:
	cd checkout && make build
	cd loms && make build
	cd notification && make build

run-all: build-all
	docker compose up -d --force-recreate --build
	cd checkout && exec ./migration.sh
	cd loms && exec ./migration.sh

reg-push: precommit build-all
	cd checkout && make docker-push
	cd loms && make docker-push

migrate:
	cd checkout && exec ./migration.sh
	cd loms && exec ./migration.sh

precommit:
	cd checkout && go mod tidy && make precommit
	cd loms && go mod tidy && make precommit
	cd notification && go mod tidy && make precommit
