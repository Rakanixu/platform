default: dependencies build

dependencies:
	go get -t -d -v ./...
build:
	cd srv && make && cd ..
	cd api && make && cd ..
protoc:
	protoc -I$$GOPATH/src --go_out=plugins=micro:$$GOPATH/src $$PWD/srv/proto/**/*.proto

