NAME=recollection

clean:
	rm -f $(NAME)

build-local:
	go build

run-local: build-local
	./$(NAME)

run-container:
	docker compose up $(NAME) --build

test:
	go test ./... -race -cover -v 2>&1