NAME=recollection

clean:
	rm -f $(NAME)

build-local:
	go build

run-local: build-local
	./$(NAME)

re-compose:
	docker compose up $(NAME) --build

compose:
	docker compose up 

test:
	go test ./... -race -cover -v 2>&1

db-shell:
	psql -h localhost -p 5432 -U postgres