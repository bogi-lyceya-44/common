package transactor

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Package transactor provides a mechanism to manage PostgreSQL transactions
// in a reusable, domain-driven way. Transactions are injected into the context
// and can be reused across layers (e.g., service and repository).
//
// Repositories should not depend on the caller to always pass a transaction.
// Instead, they should attempt to extract a transaction from the context,
// and if it's not found — begin and manage their own.
//
// Example usage in use-case/service:
//
//	err := transactor.WithTx(ctx, func(ctx context.Context) error {
//	    err := repo.CreateSomething(ctx, entity) // ctx carries tx
//	    return err
//	})
//
// In repository:
//
//	func (r *Repo) CreateSomething(ctx context.Context, e entity.Thing) error {
//	    var (
//		tx  pgx.Tx
//		err error
//		)
//
//      // if there is not tx in ctx
//      // start local tx
//		if tx, err = transactor.ExtractTx(ctx); err != nil {
//			tx, err = r.Pool.Begin(ctx)
//			if err != nil {
//				return entity.Author{}, err
//			}
//
//			defer func(tx pgx.Tx, ctx context.Context) {
//				if txErr != nil {
//					tx.Rollback(ctx)
//				}
//
//				tx.Commit(ctx)
//			}(tx, ctx)
//		}
//
//	    // main logic
//      ...
//	}
//
// This pattern allows your service layer to control transaction boundaries,
// while keeping repositories decoupled and reusable (unit-tested independently).

type Transactor interface {
	WithTx(ctx context.Context, fn func(context.Context) error) error
	WithTxOpts(ctx context.Context, fn func(context.Context) error, opts pgx.TxOptions) error
}

var _ Transactor = (*impl)(nil)

type impl struct {
	db *pgxpool.Pool
}

// Do fn in transaction with default options
func (i impl) WithTx(ctx context.Context, fn func(context.Context) error) error {
	return i.WithTxOpts(ctx, fn, pgx.TxOptions{})
}

// Do fn in transaction with custom options
func (i impl) WithTxOpts(ctx context.Context, fn func(context.Context) error, opts pgx.TxOptions) (txErr error) {
	ctxWithTx, tx, err := injectTx(ctx, i.db, opts)
	if err != nil {
		return fmt.Errorf("inject tx: %w", err)
	}

	defer func() {
		if txErr != nil {
			if rbErr := tx.Rollback(ctxWithTx); rbErr != nil {
				txErr = fmt.Errorf("tx rollback failed: %v; original error: %w", rbErr, txErr)
			}
			return
		}

		if cmErr := tx.Commit(ctxWithTx); cmErr != nil {
			txErr = fmt.Errorf("tx commit failed: %w", cmErr)
		}
	}()

	err = fn(ctxWithTx)
	if err != nil {
		return fmt.Errorf("fn: %w", err)
	}

	return nil
}

func New(db *pgxpool.Pool) *impl {
	return &impl{
		db: db,
	}
}

type txInjector struct{}

// Extracting transaction from context
// If transaction not found will return error
func ExtractTx(ctx context.Context) (pgx.Tx, error) {
	tx, ok := ctx.Value(txInjector{}).(pgx.Tx)
	if !ok {
		return nil, fmt.Errorf("tx not found in context")
	}

	return tx, nil
}

func injectTx(ctx context.Context, db *pgxpool.Pool, opts pgx.TxOptions) (context.Context, pgx.Tx, error) {
	if tx, err := ExtractTx(ctx); err == nil {
		return ctx, tx, nil
	}

	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	return context.WithValue(ctx, txInjector{}, tx), tx, nil
}
