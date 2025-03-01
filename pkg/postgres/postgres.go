package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/almostinf/glow-reminder/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

type (
	// Postgres represents a Postgres instance with pool and transaction manager.
	Postgres struct {
		// Use Pool for better performance if there are no transactions.
		Pool *pgxpool.Pool
		// Use TrManager to work with transactions.
		TrManager *manager.Manager

		logger logger.Logger
	}
)

// New creates a new Postgres instance with given url and functional options.
func New(ctx context.Context, cfg Config, logger logger.Logger) (*Postgres, error) {
	p := &Postgres{
		logger: logger,
	}

	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pool config: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.MaxPoolSize)

	for cfg.ConnAttempts > 0 {
		p.Pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
		if err == nil {
			break
		}

		p.logger.Info("Postgres is trying to connect...", map[string]interface{}{
			"attempts_left": cfg.ConnAttempts,
		})

		time.Sleep(cfg.ConnTimeout)

		cfg.ConnAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres pool: %w", err)
	}

	p.TrManager, err = manager.New(trmpgx.NewDefaultFactory(p.Pool))
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction manager: %w", err)
	}

	return p, nil
}

func (p *Postgres) GetTransactionConn(ctx context.Context) trmpgx.Tr {
	return trmpgx.DefaultCtxGetter.DefaultTrOrDB(ctx, p.Pool)
}

// Close closes the Postgres pool.
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
