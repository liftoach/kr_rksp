goose -dir ./migrations postgres "postgres://postgres:postgres@localhost:5432/finance?sslmode=disable" up 
накат
откат
goose -dir ./migrations postgres "postgres://postgres:postgres@localhost:5432/finance?sslmode=disable" down

