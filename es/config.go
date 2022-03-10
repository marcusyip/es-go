package es

import (
	_ "github.com/lib/pq"
)

type Config struct {
	TableName string
}

func NewConfig() *Config {
	return &Config{TableName: "events"}
}
