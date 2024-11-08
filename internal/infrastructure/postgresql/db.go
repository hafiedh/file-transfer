package postgresql

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/gommon/color"
	"hafiedh.com/downloader/internal/config"
)

type (
	PostgresImpl struct {
		dbConf *config.PostgreSQLDB
		Pool   *pgxpool.Pool
	}
)

func NewConnection(dbConf *config.PostgreSQLDB) (PostgresImpl, error) {
	pool, err := connect(dbConf)
	if err != nil {
		slog.Error("Failed to create connection pool", "error", err)
		return PostgresImpl{}, err
	}

	return PostgresImpl{
		dbConf: dbConf,
		Pool:   pool,
	}, nil
}

func connect(config *config.PostgreSQLDB) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.Username, config.Password, config.Name)

	dbConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		slog.Error("Failed to parse config", "error", err)
		panic(err)
	}

	dbConfig.MaxConns = int32(config.DefaultMaxConn)
	dbConfig.MinConns = int32(config.DefaultMinConn)
	dbConfig.MaxConnLifetime = config.DefaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = config.DefaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = config.DefaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = config.DefaultConnectTimeout
	dbConfig.BeforeClose = func(conn *pgx.Conn) {
		slog.Info("Closing connection to database")
	}
	dbConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		slog.Info("Connected to database")
		return nil
	}

	pool, err = pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		slog.Error("Failed to create connection pool", "error", err)
		panic(err)
	}

	// check connection
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		slog.Error("Failed to acquire connection", "error", err)
		panic(err)
	}

	defer conn.Release()

	color.Println(color.Green(fmt.Sprintf("⇨ PostgreSQL connected to database %s", config.Name)))

	return pool, nil
}

func (p *PostgresImpl) Close() error {
	p.Pool.Close()
	return nil
}
