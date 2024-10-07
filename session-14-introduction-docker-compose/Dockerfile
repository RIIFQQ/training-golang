# stage 1 : build stage
FROM golang:1.22.4-alpine AS build

# Set destination for COPY
WORKDIR /app

#copy resource code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o ./go-hello-world .

# stage 2 : run stage
FROM alpine:edge

# Set the working directory
WORKDIR /app

# Copy the binary from build stage to run stage
COPY --from=build /app/go-hello-world .

#optional : expose port
EXPOSE 8080

#run
CMD ["./go-hello-world"]