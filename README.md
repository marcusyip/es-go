![build](https://github.com/marcusyip/es-go/actions/workflows/go.yml/badge.svg)

## Overview

es-go is an event sourcing library with simplied CQRS implementation
- using Postgresql
- `One database atomic transaction` to store event and projection

## Problems to solve

- Current CQRS implementation mostly handle projection in another DB
  - event and projection are near real time consistence only
  - some use cases need strong consistence

### Data Aspects

| Aspects | Examples | 
| ---- | -------- | 
| Diff based , Latest State | - Git Commit (Diff based)<br> - Delta (Diff based - Rich Text Editing)<br> - Database WAL (Diff based) |
| Data knowledge | Postgresql does not have knowledge on JSON or JSONB columns |

### es-go libraray - positioning on data
- treat the event data (Diff based data) as the source of truth. Transform and aggregate events to projection (Latest State data)
- Application System has the knowledge on data schema but database not


## How to run test

### Create table in local DB

```
export POSTGRESQL_URL='postgres://postgres:postgres@localhost:5432/es_go_local?sslmode=disable'

migrate -database ${POSTGRESQL_URL} -path db/migrations up
```

## Remarks

1. Mircosoft - apply-simplified-microservice-cqrs-ddd-patterns
  - https://docs.microsoft.com/en-us/dotnet/architecture/microservices/microservice-ddd-cqrs-patterns/apply-simplified-microservice-cqrs-ddd-patterns
