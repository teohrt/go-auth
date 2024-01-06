# use official Golang image
FROM golang:1.21.5-alpine3.18

# set working directory
WORKDIR /app

# Copy the source code
COPY . . 

# Download and install the dependencies
RUN go get -d -v ./...

# Build the Go app
RUN go build -o recollection .

#EXPOSE the port
EXPOSE 80

# Run the executable
CMD ["./recollection"]