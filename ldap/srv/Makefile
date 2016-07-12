all: clean build

clean:
	rm -rf srv
dep:
	go get -v -d -t ./..
protoc:
	protoc -I$$GOPATH/src --go_out=plugins=micro:$$GOPATH/src $$PWD/proto/ldap.proto
build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo  .
	docker build -t kazoup/ldap-srv .
