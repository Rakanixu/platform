default: clean dep  build

clean:
	rm -rf analytics

dep: 
	go get -t -d -v  ./...
build:
	mkdir -p frontend/dist/sections/analytics
	cp -R -f ../frontend/dist/sections/analytics frontend/dist/sections
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo  .
	docker build -t kazoup/analytics-web .
