default: clean  build docker

clean:
	rm -rf api

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo  .
docker:
	docker build -t kazoup/flag-api .
