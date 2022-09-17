// Package postgres implements postgres connection.
package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	pgx_5 "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/art-injener/iot-platform/internal/config"
	"github.com/art-injener/iot-platform/pkg/logger"
)

const (
	defaultMaxPoolSize = 1
)

type Postgres struct {
	Pool        *pgxpool.Pool
	execTimeout time.Duration
	Logger      *logger.Logger
}

func NewClient(ctx context.Context, cfg *config.Config) (*Postgres, error) {
	poolConfig := prepareConfig(cfg.DBConfig)

	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		msg := fmt.Sprintf("connest to postgres %q: %v", poolConfig.ConnString(), err)
		return nil, errors.New(msg)
	}

	return &Postgres{
		Pool:        pool,
		execTimeout: time.Duration(cfg.DBConfig.ExecTimeout) * time.Second,
		Logger:      cfg.Log,
	}, nil
}

func (p *Postgres) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	ctxwt := p.getTimeoutContext(ctx)
	command, err := p.Pool.Exec(ctxwt, sql, args...)
	if err != nil {
		p.Logger.Error().Msg(fmt.Sprintf("executing sql request error: %v", err))
		return command, err
	}
	return command, nil
}

func (p *Postgres) Query(ctx context.Context, sql string, args ...interface{}) (pgx_5.Rows, error) {
	ctxwt := p.getTimeoutContext(ctx)
	rows, err := p.Pool.Query(ctxwt, sql, args...)
	if err != nil {
		p.Logger.Error().Msg(fmt.Sprintf("query sql process error: %v", err))
		return nil, err
	}
	return rows, nil
}

func (p *Postgres) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx_5.Row {
	ctxwt := p.getTimeoutContext(ctx)
	return p.Pool.QueryRow(ctxwt, sql, args...)
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

func (p *Postgres) getTimeoutContext(ctx context.Context) context.Context {
	ctxwt, _ := context.WithTimeout(ctx, p.execTimeout)
	return ctxwt
}

func prepareConfig(cfg *config.DBConfig) *pgxpool.Config {
	pgxConf, _ := pgxpool.ParseConfig("")
	pgxConf.ConnConfig.Host = cfg.Host
	pgxConf.ConnConfig.Port = cfg.Port
	pgxConf.ConnConfig.User = cfg.User
	pgxConf.ConnConfig.Database = cfg.NameDB
	pgxConf.ConnConfig.Password = cfg.Password
	pgxConf.ConnConfig.Password = cfg.Password
	pgxConf.MaxConns = defaultMaxPoolSize

	pgxConf.BeforeAcquire = func(ctx context.Context, conn *pgx_5.Conn) bool {
		return !conn.IsClosed()
	}

	return pgxConf
}
