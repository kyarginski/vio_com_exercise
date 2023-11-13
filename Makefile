GO111MODULE := auto
export GO111MODULE

lint:
	golangci-lint run ./...

test:
	go test -count=1 -race -cover ./...

test_integration:
	docker-compose up vio-db -d && make test && docker-compose down

build_loader:
	go build -tags musl -ldflags="-w -extldflags '-static' -X 'main.Version=$(VERSION)'" -o loader vio/cmd/loader

build_geolocation:
	go build -tags musl -ldflags="-w -extldflags '-static' -X 'main.Version=$(VERSION)'" -o geolocation vio/cmd/geolocation

check-swagger:
	which swagger

swagger: check-swagger
	GO111MODULE=on go mod vendor && GO111MODULE=off swagger generate spec -o ./doc/swagger.json --scan-models

serve-swagger: check-swagger
	swagger serve -F=swagger ./doc/swagger.json
