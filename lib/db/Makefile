default: protoc

protoc:
	protoc -I$$GOPATH/src --go_out=plugins=micro:$$GOPATH/src $$PWD/operations/proto/operations/*.proto
	protoc -I$$GOPATH/src --go_out=plugins=micro:$$GOPATH/src $$PWD/config/proto/config/*.proto

