default: dep vet test build

dep:
	go get -d  -v -t ./...
test:
	go test -v ./...
vet:
	go vet -v ./...
build:
	cd smtp && make && cd ..
