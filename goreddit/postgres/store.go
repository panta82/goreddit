package postgres

import (
	"context"
	"fmt"
	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	*ThreadStore
	*PostStore
	*CommentStore
}

type PgLogger struct{}

func (this *PgLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	info := ""
	if data != nil && data["sql"] != nil {
		info = fmt.Sprintf(": %s", data["sql"])
	}
	fmt.Printf("%s: %s%s\n", level, msg, info)
}

func NewStore(dataSourceName string) (*Store, error) {
	connConfig, err := pgx.ParseConfig(dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("Invalid PostgreSQL connection string \"%s\": %v", dataSourceName, err)
	}
	connConfig.Logger = &PgLogger{}
	connConfig.LogLevel = pgx.LogLevelTrace

	pgConn := stdlib.OpenDB(*connConfig)

	db := sqlx.NewDb(pgConn, "postgres")
	if err != nil {
		return nil, fmt.Errorf("Failed to open database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Failed to ping database: %v", err)
	}

	return &Store{
		ThreadStore:  &ThreadStore{db},
		PostStore:    &PostStore{db},
		CommentStore: &CommentStore{db},
	}, nil
}
