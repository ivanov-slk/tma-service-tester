FROM golang:latest

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /run

# Run
CMD ["/run/tma-service-tester"]
