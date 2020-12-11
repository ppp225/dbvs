# how to run

```
docker-compose up
go run main.go
```

## curl test

```
curl -X POST http://localhost:8090/v1/items -H "Content-type: application/json" -d '{ "name": "lorem ipsum item", "description": "is existing"}'
curl http://localhost:8090/v1/items
curl http://localhost:8090/v1/items/0x1
```
