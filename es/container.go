package es

import (
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

type Container struct {
	config              *Config
	db                  *pgxpool.Pool
	aggregateRepository AggregateRepository[AggregateRoot]
}

func NewContainer(config *Config, db *pgxpool.Pool, aggregateRepository AggregateRepository[AggregateRoot]) *Container {
	return &Container{
		config:              config,
		db:                  db,
		aggregateRepository: aggregateRepository,
	}
}

func (c *Container) GetConfig() *Config {
	return c.config
}

func (c *Container) GetDB() *pgxpool.Pool {
	return c.db
}

func (c *Container) GetAggregateRepository() AggregateRepository[AggregateRoot] {
	return c.aggregateRepository
}
