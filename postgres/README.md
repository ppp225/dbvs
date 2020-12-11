## machine environment

- migrate-cli
https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
```sh
go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate
```

## create new migration

```
migrate create -ext sql -dir db/migrations -seq create_items_table
```

## run migrations on db

```
export POSTGRES_URL="postgres://user:pass@localhost:5432/test_db?sslmode=disable"
migrate -database ${POSTGRES_URL} -path db/migrations up
```

# how to run

```
docker-compose up
#<migrations>
go run main.go
```

## curl test

```
curl -X POST http://localhost:8090/v1/items -H "Content-type: application/json" -d '{ "name": "lorem ipsum item", "description": "is existing"}'
curl http://localhost:8090/v1/items
curl http://localhost:8090/v1/items/1
```
