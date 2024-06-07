package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ravenocx/hospital-mgt/config"
)

func OpenConnection(config config.Config) (*pgxpool.Pool, error) {
	ctx := context.Background()

	postgresUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.DbUsername,
		config.DbPassword,
		config.DbHost,
		config.DbPort,
		config.DbName,
	)

	dbconfig, err := pgxpool.ParseConfig(postgresUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to parse pool config : %+v", err)
	}

	dbconfig.MaxConnLifetime = time.Duration(config.DbMaxLifetimeConn) * time.Hour
	dbconfig.MaxConnIdleTime = time.Duration(config.DbMaxIdleConn) * time.Minute
	dbconfig.HealthCheckPeriod = 5 * time.Second
	dbconfig.MaxConns = int32(config.DbMaxConn)
	dbconfig.MinConns = 10

	pool, err := pgxpool.NewWithConfig(ctx,dbconfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create db pool : %+v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database : %+v", err)
	}

	return pool, nil
}