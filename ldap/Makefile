default: dep build

dep: 
	go get -v -t -d ./...
build:
	cd srv && make && cd ..
protoc:
	protoc -I$$GOPATH/src --go_out=plugins=micro:$$GOPATH/src $$PWD/srv/proto/**/*.proto

