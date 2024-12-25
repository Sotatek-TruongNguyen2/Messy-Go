package db

import (
	db_utils "arthur/simple_bank/db/utils"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)


func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams {
		Owner: db_utils.RandomOwner(),
		Balance: db_utils.RandomMoney(),
		Currency: db_utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	
	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	newAccount := createRandomAccount(t)
	account, err := testQueries.GetAccount(context.Background(), newAccount.ID)

	if (err != nil) {
		t.Fatal("Error when trying to get account details by id")
	}

	require.Equal(t, newAccount.Balance, account.Balance)
	require.Equal(t, newAccount.Owner, account.Owner)
}

// func TestBackground(t *testing.T) {
// 	ctx := context.Background()
// 	pointer := unsafe.Pointer(&ctx)

// 	underlyingPointer := *(*int)(pointer)

// 	t.Log("Rafael :", &ctx)
// 	t.Logf("%p %v %p", pointer, &pointer, (*int)(pointer))
// 	t.Logf("%v 0x%x %v", underlyingPointer, underlyingPointer, &underlyingPointer)
// 	halo(ctx, t)

// 	t.Logf("Rafael cooked: %v", *(*int)(pointer))
// }


// func halo(ctx context.Context, t *testing.T) {
// 	pointer := unsafe.Pointer(&ctx)
// 	underlyingPointer := *(*int)(pointer)

// 	t.Log("Rafael V2:", &ctx)
// 	t.Logf("%v 0x%x %v", underlyingPointer, underlyingPointer, &underlyingPointer)

// 	t.Logf("%v %p", *(*int)(pointer), (*int)(pointer))
// }