default: clean dep build

clean:
	rm -rf search

dep: 
	go get -t -d -v  ./...
build:
	mkdir -p frontend/dist/sections/search
	cp -R -f ../frontend/dist/sections/search frontend/dist/sections
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo  .
	docker build -t kazoup/search-web .
