.PHONY: default

default: di build;

di:
	wire ./cmd/cli
	wire ./cmd/airflow

build:
	go build -o out/cli ./cmd/cli
	go build -o out/airflow ./cmd/airflow

install:
	go install ./cmd/cli
	go install ./cmd/airflow

run:
	go run ./cmd/cli
	go run ./cmd/airflow

static:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./out/cli ./cmd/cli
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./out/airflow ./cmd/airflow
