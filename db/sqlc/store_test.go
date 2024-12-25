package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// run n concurrent transfer transactions
	n := 10
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i % 2 != 0 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		txName := fmt.Sprintf("tx %d", i + 1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID: toAccountID,
				Amount: amount,
			})

			errs <- err 
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		select {
			case err := <- errs: 
				require.NoError(t, err)
			case result := <- results:
				require.NotEmpty(t, result)
				
				transfer := result.Transfer

				require.NotEmpty(t, transfer)
				// require.Equal(t, account1.ID, transfer.FromAccountID.Int64)
				// require.Equal(t, account2.ID, transfer.ToAccountID.Int64)
				// require.Equal(t, amount, transfer.Amount)

				// require.NotZero(t, transfer.ID)
				// require.NotZero(t, transfer.CreatedAt)

				// _, err := store.GetTransfer(context.Background(), transfer.ID)
				// require.NoError(t, err)


				// fromAccount := result.FromAccount
				// require.NotEmpty(t, fromAccount)
				// require.Equal(t, account1.ID, fromAccount.ID)


				// toAccount := result.ToAccount
				// require.NotEmpty(t, toAccount)
				// require.Equal(t, account2.ID, toAccount.ID)

				// diff1 := account1.Balance - fromAccount.Balance
				// diff2 := toAccount.Balance - account2.Balance

				// require.Equal(t, diff1, diff2)
		}
	}
}