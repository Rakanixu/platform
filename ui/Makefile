default: build


build-front:
	cd frontend && gulp && cd ..
	cd frontend && npm install && bower install && cd ..
clean:
	rm -rf frontend dist
run:
	go run  main.go  --registry=mdns
dist:
	cd frontend && npm install && bower install && cd ..

build:
	cd frontend && npm install && npm install gulp && bower install && gulp && cd ..
	cd static && make && cd ..
	cd dashboard && make && cd ..
	cd analytics && make && cd ..
	cd search && make && cd ..
	cd logs && make && cd ..
	cd docs && make && cd ..
	cd login && make && cd ..
	cd settings && make && cd ..
	cd wizard && make && cd ..
	cd static && rm -rf frontend && cd ..
	cd dashboard && rm -rf frontend && cd ..
	cd analytics && rm -rf frontend && cd ..
	cd search && rm -rf frontend && cd ..
	cd logs && rm -rf frontend && cd ..
	cd docs && rm -rf frontend && cd ..
	cd login && rm -rf frontend && cd ..
	cd settings && rm -rf frontend && cd ..
	cd wizard && rm -rf frontend && cd ..

