default: clean dep  build

clean:
	rm -rf docs

dep: 
	go get -t -d -v  ./...
build:
	mkdir -p frontend/dist/sections/docs
	cp -R -f ../frontend/dist/sections/docs frontend/dist/sections
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo  .
	docker build -t kazoup/docs-web .
