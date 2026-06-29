package database

import (
	"fmt"
	"time"
)

type Config struct {
	DSN          string
	MaxIdle      int
	MaxOpen      int
	MaxLifeTime  time.Duration
	MaxIdleTime  time.Duration
	MaxBusy      int
	QueryTimeout time.Duration
}

func NewConfigWithDSN(dsn string) *Config {
	return &Config{
		DSN:          dsn,
		MaxIdle:      10, // In a WAL environment MaxIdle must == MaxOpen,
		MaxOpen:      10, // otherwise, the extra connections are torn down when demand drops
		MaxLifeTime:  60 * time.Second,
		MaxIdleTime:  30 * time.Second,
		MaxBusy:      5000,
		QueryTimeout: 5 * time.Second,
	}
}

func (cfg *Config) IsValid() error {
	if cfg.DSN == "" {
		return fmt.Errorf("database config: missing DSN")
	}

	if cfg.MaxIdle == 0 {
		return fmt.Errorf("database config: MaxIdle must be a positive integer")
	}

	if cfg.MaxOpen == 0 {
		return fmt.Errorf("database config: MaxOpen must be a positive integer")
	}

	if cfg.MaxLifeTime == 0 {
		return fmt.Errorf("database config: MaxLifeTime must be a valid time.Duration")
	}

	if cfg.MaxIdleTime == 0 {
		return fmt.Errorf("database config: MaxIdleTime must be a valid time.Duration")
	}

	if cfg.MaxBusy == 0 {
		return fmt.Errorf("database config: MaxBusy must be a positive integer")
	}

	if cfg.QueryTimeout == 0 {
		return fmt.Errorf("database config: QueryTimeout must be a valid time.Duration")
	}

	return nil
}
