.PHONY: build
build:
	CGO_ENABLED=0 go build -o slowdns .

.PHONY: docker-build
docker-build: GOOS = linux
docker-build: GARCH = amd64 
docker-build: build
	docker build -t kopiczko/slowdns .
