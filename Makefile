build:
	docker build . -t docker-example:latest

run:
	docker run -p 8080:8080 docker-example:latest

test:
	go test -v ./...