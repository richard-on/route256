goose -dir ./migrations postgres "postgres://admin:pgpswd@localhost:6442/loms?sslmode=disable" status

goose -dir ./migrations postgres "postgres://admin:pgpswd@localhost:6442/loms?sslmode=disable" up