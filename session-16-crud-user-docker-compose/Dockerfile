# stage 1 : build stage
FROM golang:1.22.4-alpine AS build

#install necessary packages
RUN apk add --no-cache git make protobuf


# Set destination for COPY
WORKDIR /app

#copy the go module manifest and download the dependencies
COPY go.mod go.sum ./
RUN go mod download

#copy resource code
COPY . .

#intsall golangci-lint
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1

#run the linter with increased timeout
RUN golangci-lint run --timeout 5m

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o ./go-crud-user .

# stage 2 : run stage
FROM alpine:edge

# Set the working directory
WORKDIR /app

# Copy the binary from build stage to run stage
COPY --from=build /app/go-crud-user .

#install CA certificates
RUN apk --no-cache add ca-certificates

#run
CMD ["./go-crud-user"]