BINARY=engine
test: 
	go test -v -cover -covermode=atomic ./...

engine:
	go build -o ${BINARY} main.go

unittest:
	go test -short  ./...

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

compile:
	go mod tidy
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/app

docker:
	docker build . -t app:1.0.0

run:
	docker run -p 3000:3000 -d app:1.0.0

stop:
	docker-compose down

lint-prepare:
	@echo "Installing golangci-lint" 
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

lint:
	./bin/golangci-lint run \
		--exclude-use-default=false \
		--enable=golint \
		--enable=gocyclo \
		--enable=goconst \
		--enable=unconvert \
		./...

.PHONY: clean install unittest build docker run stop vendor lint-prepare lint