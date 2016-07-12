default: clean deps test build docker

clean:
	rm -rf web-proxy
deps:
	go get -d -v ./...
test: 
	go test  ./...
build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo  .
docker:
	docker build -t kazoup/web-proxy .
deploy:
	docker push kazoup/web-proxy
