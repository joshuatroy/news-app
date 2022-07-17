# news-app
Backend Engineer Tech Test for building a news app API to power an app using RSS feeds

# Contributing

To install dependencies:
```go
go mod vendor
```

To generate mocks:
```go
go generate ./...
```

To run tests:
```go
go test ./...
```

To run the application:
```go
go run ./cmd/server/main.go
```

To compile and execute the build:
```go
go build ./cmd/server
./server
```

To test endpoints using postman please import **postman_collection.json** file