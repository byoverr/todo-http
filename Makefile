build:
	go build -o bin/main cmd/todo-server/main.go

run:
	go run cmd/todo-server/main.go

docker-build-and-run:
	docker build -t todo-server .
	docker run --name todo-app -p 8080:8080 todo-server

lint:
	golangci-lint run
test:
	go test -v ./...