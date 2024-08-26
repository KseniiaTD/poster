package database

import (
	"fmt"
	"os"

	"github.com/KseniiaTD/poster/internal/database/common"
	"github.com/KseniiaTD/poster/internal/database/inmemory"
	"github.com/KseniiaTD/poster/internal/database/postgres"
)

func New(isInMemory bool) (common.Database, error) {
	if !isInMemory {
		dsn := getDSN()

		return postgres.New(dsn)
	}

	return inmemory.New(), nil
}

func getDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))
}
