package database

import (
	"fmt"

	"github.com/KseniiaTD/poster/config"
	"github.com/KseniiaTD/poster/internal/database/common"
	"github.com/KseniiaTD/poster/internal/database/inmemory"
	"github.com/KseniiaTD/poster/internal/database/postgres"
)

func New(isInMemory bool, cfg config.Config) (common.Database, error) {
	if !isInMemory {
		dsn := getDSN(cfg)

		return postgres.New(dsn)
	}

	return inmemory.New(), nil
}

func getDSN(cfg config.Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.User,
		cfg.Password,
		cfg.DB)
}
