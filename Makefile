BINARY_NAME=recollection
LOCAL_SERVER_PORT=8000
BASE_ENV_VALS := SERVER_PORT="$(LOCAL_SERVER_PORT)"

build:
	go build

clean:
	rm -f $(BINARY_NAME)

run: build
	$(BASE_ENV_VALS) ./$(BINARY_NAME)

test:
	go test ./... -race -cover -v 2>&1