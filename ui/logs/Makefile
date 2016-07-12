default: clean  build

clean:
	rm -rf logs

build:
	mkdir -p frontend/dist/sections/logs
	cp -R -f ../frontend/dist/sections/logs frontend/dist/sections
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo  .
	docker build -t kazoup/logs-web .
