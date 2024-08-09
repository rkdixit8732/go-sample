# Start with a base image that includes Go
FROM golang:latest

# Set the working directory inside the container
WORKDIR /go/src/app

# Copy the local package files to the container's workspace
COPY . .

# Fetch dependencies of golang
RUN go get -d -v ./...

# Build the Go application
RUN go build -o myapp

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./myapp"]