default: clean dep  build

clean:
	rm -rf static
dep:
	go get -v -d -t ./...
build:
	mkdir -p frontend
	cp -R -f ../frontend/dist frontend
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo  .
	docker build -t kazoup/static-web .
