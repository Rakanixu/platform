default: clean build docker

clean:
	rm -rf srv
dependencies:
	go get -v ./...
build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo  .
docker:
	docker build -t kazoup/flag-srv .
protoc:
	protoc -I$$GOPATH/src --go_out=plugins=micro:$$GOPATH/src $$PWD/proto/**/*.proto

