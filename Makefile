run:
	docker-compose up -d
	
run-producer:
	go run ./cmd/producer/main.go
