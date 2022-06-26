package es

import (
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Container struct {
	config              *Config
	logger              *zap.Logger
	db                  *pgxpool.Pool
	aggregateRepository AggregateRepository[AggregateRoot]
}

func NewContainer(config *Config, logger *zap.Logger, db *pgxpool.Pool,
	aggregateRepository AggregateRepository[AggregateRoot],
) *Container {
	return &Container{
		config:              config,
		logger:              logger,
		db:                  db,
		aggregateRepository: aggregateRepository,
	}
}

func (c *Container) GetConfig() *Config {
	return c.config
}

func (c *Container) GetLogger() *zap.Logger {
	return c.logger
}

func (c *Container) GetDB() *pgxpool.Pool {
	return c.db
}

func (c *Container) GetAggregateRepository() AggregateRepository[AggregateRoot] {
	return c.aggregateRepository
}
