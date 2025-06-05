run:
	go run cmd/main.go

test:
	go test ./...

swag:
	swag init  -g cmd/main.go -o cmd/docs/
