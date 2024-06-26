.PHONY: start
start:
	./service.exe -config service.config.json start

.PHONY: stop
stop:
	./service.exe -config service.config.json stop

.PHONY: install
install:
	./service.exe -config service.config.json install

.PHONY: delete
delete:
	./service.exe -config service.config.json delete

.PHONY: build
build:
	GOOS=windows go build -o service.exe cmd/main.go
	GOOS=windows go build -o server.exe test/test_server/cmd/main.go

.PHONY: update
update:
	go mod tidy
	go mod vendor
	go fmt ./...

.PHONY: test
test:
	GOOS=windows go build -o test/test_service/cmd/test_service.exe test/test_service/cmd/main.go
	GOOS=windows go build -o test/test_server/cmd/test_server.exe test/test_server/cmd/main.go
	go test -count=1 -v ./test/...
