![build](https://github.com/marcusyip/es-go/actions/workflows/go.yml/badge.svg)

## Create table in local DB
export POSTGRESQL_URL='postgres://postgres:postgres@localhost:5432/es_go_local?sslmode=disable'                     Go 1.18

migrate -database ${POSTGRESQL_URL} -path db/migrations up
