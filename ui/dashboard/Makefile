default: clean dep  build

clean:
	rm -rf dashboard

dep: 
	go get -t -d -v  ./...
build:
	mkdir -p frontend/dist/sections/dashboard
	cp -R -f ../frontend/dist/sections/dashboard frontend/dist/sections
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo  .
	docker build -t kazoup/dashboard-web .
