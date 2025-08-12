package transactor

import (
	"context"
	stdErrors "errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	pkgErrors "github.com/pkg/errors"
)

// Tx key in context
type txInjector struct{}

type Querier interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Transactor struct {
	db *pgxpool.Pool
}

// Do fn in transaction with default options
func (i Transactor) WithTx(
	ctx context.Context,
	fn func(context.Context) error,
) error {
	return i.WithTxOpts(ctx, fn, pgx.TxOptions{})
}

// Do fn in transaction with custom options
func (i Transactor) WithTxOpts(
	ctx context.Context,
	fn func(context.Context) error,
	opts pgx.TxOptions,
) (txErr error) {
	ctxWithTx, tx, err := injectTx(ctx, i.db, opts)
	if err != nil {
		return pkgErrors.Wrap(err, "inject tx")
	}

	defer func() {
		if txErr != nil {
			txErr = stdErrors.Join(
				txErr,
				pkgErrors.Wrap(tx.Rollback(ctxWithTx), "rollback tx"),
			)
		}
	}()

	err = fn(ctxWithTx)
	if err != nil {
		return pkgErrors.Wrap(err, "tx function")
	}

	return pkgErrors.Wrap(tx.Commit(ctxWithTx), "commit tx")
}

// Extracting transaction from context
// If transaction not found will return Pool
// ExtractTx returns the query interface
func (t Transactor) ExtractTx(ctx context.Context) Querier {
	tx, ok := ctx.Value(txInjector{}).(pgx.Tx)
	if !ok {
		return t.db
	}

	return tx
}

func New(db *pgxpool.Pool) *Transactor {
	return &Transactor{
		db: db,
	}
}

// Injecting transaction into context
// If transaction is alreay in context, don't create a new one
func injectTx(
	ctx context.Context,
	db *pgxpool.Pool,
	opts pgx.TxOptions,
) (context.Context, pgx.Tx, error) {
	// try to extract transaction if exists
	if tx, ok := ctx.Value(txInjector{}).(pgx.Tx); ok {
		return ctx, tx, nil
	}

	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return nil, nil, pkgErrors.Wrap(err, "begin tx")
	}

	return context.WithValue(ctx, txInjector{}, tx), tx, nil
}
