CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
SMARTIMPORTS=${BINDIR}/smartimports_${GOVER}
LINTVER=v1.51.1
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}
PACKAGE=gitlab.ozon.dev/rragusskiy/homework-1/loms/cmd/app
RPCNAME=loms.proto
PROTOPATH=api/loms/v1
RPCPATH=pkg/loms
PROTOVER=v1
BUILDDATE=$(shell date -u +%Y%m%d.%H%M%S)
APPVER=0.6.1

all: format build test lint

build: bindir
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -ldflags="-X gitlab.ozon.dev/rragusskiy/homework-1/loms/config.Version=${APPVER} \
	 -X gitlab.ozon.dev/rragusskiy/homework-1/loms/config.Build=${BUILDDATE} \
	 -X gitlab.ozon.dev/rragusskiy/homework-1/loms/config.ProtoVersion=${PROTOVER}" \
 					   		-o ${BINDIR}/app ${PACKAGE}

test:
	go test ./...

run:
	go run ${PACKAGE}

lint: install-lint
	${LINTBIN} run

precommit: format build test lint
	echo "OK"

bindir:
	mkdir -p ${BINDIR}

format: install-smartimports
	${SMARTIMPORTS} -exclude internal/mocks

test-coverage:
	@go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out; \
	go tool cover -func=./coverage.out | grep "total";

install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})

install-smartimports: bindir
	test -f ${SMARTIMPORTS} || \
		(GOBIN=${BINDIR} go install github.com/pav5000/smartimports/cmd/smartimports@latest && \
		mv ${BINDIR}/smartimports ${SMARTIMPORTS})

proto-all: install-go-deps get-go-deps vendor-proto generate

install-go-deps:
	GOBIN=$(BINDIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(BINDIR) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

get-go-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

vendor-proto:
		mkdir -p vendor-proto
		@if [ ! -d vendor-proto/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor-proto/googleapis &&\
			mkdir -p  vendor-proto/google/ &&\
			mv vendor-proto/googleapis/google/api vendor-proto/google &&\
			rm -rf vendor-proto/googleapis ;\
		fi
		@if [ ! -d vendor-proto/google/protobuf ]; then\
			git clone https://github.com/protocolbuffers/protobuf vendor-proto/protobuf &&\
			mkdir -p  vendor-proto/google/protobuf &&\
			mv vendor-proto/protobuf/src/google/protobuf/*.proto vendor-proto/google/protobuf &&\
			rm -rf vendor-proto/protobuf ;\
		fi

generate:
	mkdir -p $(RPCPATH)
	protoc -I $(PROTOPATH) -I vendor-proto \
	--go_out=$(RPCPATH) --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=$(RPCPATH) --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/loms/v1/loms.proto

migrate:
	exec ./migration.sh

docker-push:
	docker tag route256-loms richardhere/route256-loms:latest
	docker push richardhere/route256-loms

config-file:
	kubectl create configmap loms -n route256 --from-file=config/config.yaml
	kubectl create configmap odyssey-loms -n route256 --from-file=config/odyssey/odyssey.conf
	kubectl create configmap pgconf-loms -n route256 --from-file=config/postgres/postgresql.conf
	kubectl create configmap pghba-loms -n route256 --from-file=config/postgres/pg_hba.conf
