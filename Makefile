include $(PWD)/.env
NAME=gogoapps

coverage:
	@go test -failfast -race -short -coverprofile=cover.out ./...
	@go tool cover -func=cover.out | awk 'END{print $$3}' | xargs echo "code coverage:"

build:
	@GOOS=linux GOARCH=amd64 go build .

test:
	@go test -timeout 1m -race $(PACKAGES)

clean:
	@go clean -i -x

docker-build:
	docker build -t gogoapps:latest .
