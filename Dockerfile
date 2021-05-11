FROM golang:1.16.0

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on

# Setup working directory
WORKDIR /app

# Copy and download dependencies
COPY go.mod .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main .

# Export port
EXPOSE 3000

# Command to run when starting the container
CMD ["/app/main"]
