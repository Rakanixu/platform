default: dep build

dep:
	go get -u -v -t ./...
build:
	cd smtp && make && cd ..
