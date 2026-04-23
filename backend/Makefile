.PHONY: build seed-instruments seed-instruments-force
build:
	go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o build/app main.go

run:
	CONFIG_PATH=./config/dev.yaml ./build/app

run-dev:
	make build && make run

install-tools:
	go get -u github.com/swaggo/swag/cmd/swag
	go install github.com/golang/mock/mockgen@v1.6.0
	curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | apt-key add -
	echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list
	apt-get update
	apt-get install -y migrate

swagger:
	swag init --parseDependency -g main.go --output=./docs

docker-build:
	cd .. && docker compose up -d --build

migrate-up:
	CONFIG_PATH=./config/dev.yaml go run migrator.go up

migrate-down:
	CONFIG_PATH=./config/dev.yaml go run migrator.go down

# Fills the catalog with demo instruments if the table is empty. Use seed-instruments-force to append another batch.
seed-instruments:
	CONFIG_PATH=./config/dev.yaml go run ./cmd/seed-instruments

seed-instruments-force:
	CONFIG_PATH=./config/dev.yaml go run ./cmd/seed-instruments -force