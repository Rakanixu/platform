default: clean dep  build

clean:
	rm -rf web
dep: 
	go get -t -d -v  ./...

build:
	polymer build
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo  .
	docker build -t kazoup/webapp-web .
