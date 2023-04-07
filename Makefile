
test-coverage:
	go test -v ./checkout/... ./loms/... -coverprofile=coverage.out
	go tool cover -func=./coverage.out | grep "total";
	go tool cover -html=coverage.out

build-all:
	cd checkout && make build
	cd loms && make build
	cd notification && make build

.PHONY: logs
logs:
	mkdir -p logs/data
	touch logs/data/log.txt
	touch logs/data/offsets.yaml
	sudo chmod -R 777 logs/data
	cd logs && docker compose up -d

.PHONY: metrics
metrics:
	mkdir -p metrics/data
	sudo chmod -R 777 metrics/data
	cd metrics && sudo docker compose up -d

clean:
	sudo rm -rf ./logs/data/*
	sudo rm -rf ./metrics/data/*
	sudo docker volume rm prometheus


run-all: build-all
	docker compose up -d --force-recreate --build
	cd checkout && exec ./migration.sh
	cd loms && exec ./migration.sh
	# docker compose logs -f checkout loms notification | cut -f2 -d '|' > logs/data/log.txt &

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
