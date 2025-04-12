package testutil

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxtest"
)

func ConnTestRunner() pgxtest.ConnTestRunner {
	testRunner := pgxtest.DefaultConnTestRunner()

	if pgConnString := os.Getenv("TEST_PG_CONN_STRING"); pgConnString != "" {
		testRunner.CreateConfig = func(ctx context.Context, t testing.TB) *pgx.ConnConfig {
			config, err := pgx.ParseConfig(pgConnString)
			if err != nil {
				t.Fatalf("ParseConfig failed: %v", err)
			}
			return config
		}
	}

	return testRunner
}
