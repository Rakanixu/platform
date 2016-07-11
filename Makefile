default: dep build

dep:
	go get -d  -v -t ./...
build:
	cd smtp && make && cd ..
