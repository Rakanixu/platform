default: dep build

dep:
	go get -v -t ./...
build:
	cd smtp && make && cd ..
