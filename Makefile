clean:
	rm -r data/*

start:
	docker-compose up -d

stop:
	docker-compose down

init:
	docker exec stocks_db psql -U postgres -f /mnt/init.sql

run:
	go run src/main.go
