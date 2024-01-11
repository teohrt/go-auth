NAME=recollection

clean:
	rm -f $(NAME)

build-local:
	go build

run-local: build-local
	./$(NAME)

compose:
	docker compose up $(NAME) --build

test:
	go test ./... -race -cover -v 2>&1