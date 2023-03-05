goose -dir ./migration postgres "postgres://admin:pgpswd@localhost:5422/checkout?sslmode=disable" status

goose -dir ./migration postgres "postgres://admin:pgpswd@localhost:5422/checkout?sslmode=disable" up