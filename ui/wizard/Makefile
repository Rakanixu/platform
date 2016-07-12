default: clean dep build

clean:
	rm -rf wizard

dep: 
	go get -t -d -v  ./...
build:
	mkdir -p frontend/dist/sections/wizard
	cp -R -f ../frontend/dist/sections/wizard frontend/dist/sections
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo  .
	docker build -t kazoup/wizard-web .
