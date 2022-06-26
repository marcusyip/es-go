![test](https://github.com/marcusyip/es-go/actions/workflows/test.yml/badge.svg)
![lint](https://github.com/marcusyip/es-go/actions/workflows/lint.yml/badge.svg)
[![codecov](https://codecov.io/gh/marcusyip/es-go/branch/main/graph/badge.svg?token=9EP43TF3X4)](https://codecov.io/gh/marcusyip/es-go)

`Caution: the library concept used in other language and projects. BUT this golang library is an experimental implementation`

## Overview

es-go is an event sourcing library with simplied CQRS implementation
- Using Postgresql
- `One database atomic transaction` to store event and projection
- Event data as source of truth
- Considered database performance in mind
  - table indexing
  - projection for aggregate loader (load aggregate from another table BUT not aggregating from events for better performance)

`Please check test/e2e for how it used`

## Problems to solve

- Current CQRS implementation mostly handle projection in another DB
  - event and projection are near real time consistence only
  - some use cases need strong consistence

### Different Data Aspects

| Aspects | Examples | 
| ---- | -------- | 
| Diff based , Latest State | - Git Commit (Diff based)<br> - Delta (Diff based - Rich Text Editing)<br> - Database WAL (Diff based) |
| Data knowledge | Postgresql does not have knowledge on JSON or JSONB columns |

### es-go libraray
- the event data (Diff based data) as the source of truth. Transform and aggregate events to projection (Latest State data)
- Application System has the knowledge on data schema, but database does not have knowledge

## How to run test

1. Run postgresql
```
docker-compose up
```

2. database schema migration
```
export POSTGRESQL_URL='postgres://postgres:postgres@localhost:5432/es_go_local?sslmode=disable'

migrate -database ${POSTGRESQL_URL} -path db/migrations up
```

```
go run test ./...
```

## References

1. Mircosoft - apply-simplified-microservice-cqrs-ddd-patterns
  - https://docs.microsoft.com/en-us/dotnet/architecture/microservices/microservice-ddd-cqrs-patterns/apply-simplified-microservice-cqrs-ddd-patterns
