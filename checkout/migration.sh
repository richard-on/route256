goose -dir ./migrations postgres "postgres://admin:pgpswd@localhost:6422/checkout?sslmode=disable" status

goose -dir ./migrations postgres "postgres://admin:pgpswd@localhost:6422/checkout?sslmode=disable" up