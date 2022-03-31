# go-ms-beers
micro service in Go

## Installation

  ```sh
// download dependencies
go mod tidy

// Ejecutar
go run cmd/server/main.go

// run tests
go test ./internal/app/... -v

// Coverage
go test ./internal/app/... -cover
```

Go to http://localhost:3000 to see the swagger specification