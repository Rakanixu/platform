default: clean dep  build

clean:
	rm -rf login
dep: 
	go get -t -d -v  ./...

build:
	mkdir -p frontend/dist/sections/login
	cp -R -f ../frontend/dist/sections/login frontend/dist/sections
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo  .
	docker build -t kazoup/login-web .
