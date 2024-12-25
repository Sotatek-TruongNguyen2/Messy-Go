package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries
	conn *pgxpool.Pool
}


var txKey = struct{}{}

func NewStore(conn *pgxpool.Pool) *Store {
	return &Store {
		New(conn),
		conn,
	}
}

func(s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {
		return err
	}

	q := s.WithTx(tx)
	err = fn(q)

	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return fmt.Errorf("tx err :%v, rb err: %v", err, rollbackErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 `json:"to_account_id"`
	Amount int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry  `json:"to_entry"`
}

func(s *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: pgtype.Int8{
				Int64: arg.FromAccountID,
				Valid: true,
			},
			ToAccountID: pgtype.Int8{
				Int64: arg.ToAccountID,
				Valid: true,
			},
			Amount: arg.Amount,
		})

		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 1")

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: pgtype.Int8{
				Int64: arg.FromAccountID,
				Valid: true,
			},
			Amount: -arg.Amount,
		})

		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 2")

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: pgtype.Int8{
				Int64: arg.ToAccountID,
				Valid: true,
			},
			Amount: arg.Amount,
		})

		if err != nil {
			return err
		}

		// TODO: Update balance (needs to done later on as it will involve dead-lock problem)

		if (arg.FromAccountID < arg.ToAccountID) {
			result.FromAccount, result.ToAccount, err = s.TransferMoney(ctx, q, arg.FromAccountID, arg.ToAccountID, -arg.Amount, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = s.TransferMoney(ctx, q, arg.ToAccountID, arg.FromAccountID, arg.Amount, -arg.Amount)
		}

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func (s *Store) TransferMoney(
	ctx context.Context,
	q *Queries,
	fromAccountID int64,
	toAccountID int64,
	fromAmount int64,
	toAmount int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID: fromAccountID,
		Amount: fromAmount,
	})

	if err != nil {
		return
	}

	account2, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID: toAccountID,
		Amount: toAmount,
	})

	if err != nil {
		return
	}

	return
}