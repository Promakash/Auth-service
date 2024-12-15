generate_docs:
	swag init -g cmd/main.go

build: generate_docs
	docker compose build

launch:
	docker compose up -d
build_launch: build launch

build_launch_tests: build_launch
	docker compose --profile test up -d

stop:
	docker compose down