# dbvs

dbvs is a playground for benchmarking different db against each other through a basic rest api interface.

## why

it's surprisingly hard to find benchmarks.

## what this is not

this isn't a comprehensive benchmark of databases, use cases or scalability across hundreds of servers.

## how to run

1. run selected api (as seen in subdir/README)
2. run tests as in [methodology](#methodology)

## methodology

schema can be described with the following go type, however it may differ slightly based on implementation:

```go
type Item struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"createdAt"`
}
```

tests consist of running bombardier:
```
1. bombardier -c 100 -n 10000 -l -m POST -b '{ "name": "lorem ipsum item name", "description": "is existing"}' http://localhost:8090/v1/items
2. bombardier -c 100 -n 10000 -l -m GET http://localhost:8090/v1/items/<ID>
3. bombardier -c 1 -n 100 -l -m GET http://localhost:8090/v1/items
4. bombardier -c 100 -n 10000 -l -m GET http://localhost:8090/v1/items
```
1. creating 10000 objects with 100 simultaneous connections
2. get one item 10000times/100conn
3. get list of 1000 items sorted by date descending (reverse order to insert) 100 times sequentially
4. 3 but 10000times/100conn

## results

| Test | DB              | Avg req/s     | Latency avg ms | 50% | 75% | 90% | 95% | 99% |
| ----:| --------------- | -------------:| --------------:| ---:| ---:| ---:| ---:| ---:|
| 1.   | postgres 13     |          1500 |             66 |  29 | 126 | 151 | 166 | 196 |
| 1.   | dgraph v20.07.2 |          1865 |             53 |  52 |  58 |  64 |  70 |  97 |
| 2.   | postgres 13     |          1732 |             57 |  16 | 129 | 166 | 180 | 206 |
| 2.   | dgraph v20.07.2 |          4371 |             23 |  19 |  22 |  27 |  31 | 187 |
| 3.   | postgres 13     |           144 |              7 |   7 |   7 |   8 |   8 |   8 |
| 3.   | dgraph v20.07.2 |             8 |            117 | 116 | 123 | 126 | 132 | 143 |
| 4.   | postgres 13     |           645 |             80 |  79 | 277 | 399 | 444 | 526 |
| 4.   | dgraph v20.07.2 |        failed |              - |   - |   - |   - |   - |   - |

## summary

never forget to cache expensive requests

## future work

- joins
- combined r/w workload
- multi node db
- automate benchmark suite

## note

everything may or may not change ¯\\\_(ツ)\_/¯
