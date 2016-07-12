default: clean  build

clean:
	rm -rf settings

build:
	mkdir -p frontend/dist/sections/settings
	cp -R -f ../frontend/dist/sections/settings frontend/dist/sections
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo  .
	docker build -t kazoup/settings-web .
