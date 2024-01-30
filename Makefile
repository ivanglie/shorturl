tests:
	go test -v -cover -race ./...

run:
	go run -race ./cmd/app/

dc:
	docker compose up -d