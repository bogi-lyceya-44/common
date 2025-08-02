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

func (i impl) WithTx(ctx context.Context, fn func(context.Context) error) error {
	return i.WithTxOpts(ctx, fn, pgx.TxOptions{})
}

func (i impl) WithTxOpts(ctx context.Context, fn func(context.Context) error, opts pgx.TxOptions) (txErr error) {
	ctxWithTx, err, tx := injectTx(ctx, i.db, opts)
	if err != nil {
		return fmt.Errorf("inject tx: %w", err)
	}

	defer func() {
		if txErr != nil {
			tx.Rollback(ctxWithTx)
			return
		}

		tx.Commit(ctxWithTx)
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

func ExtractTx(ctx context.Context) (pgx.Tx, error) {
	tx, ok := ctx.Value(txInjector{}).(pgx.Tx)
	if !ok {
		return nil, fmt.Errorf("tx not found in context")
	}

	return tx, nil
}

func injectTx(ctx context.Context, db *pgxpool.Pool, opts pgx.TxOptions) (context.Context, error, pgx.Tx) {
	if tx, err := ExtractTx(ctx); err == nil {
		return ctx, nil, tx
	}

	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err, nil
	}

	return context.WithValue(ctx, txInjector{}, tx), nil, tx
}
