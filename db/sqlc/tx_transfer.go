package db

import (
	"context"
	"fmt"
)

type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_accont"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var ctxKey = struct{}{}

// perform transaction, record entries, update accounts
func (store *SqlStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(ctxKey)

		// create transaction
		fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// record two entries
		fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// update two account balance
		if arg.FromAccountId < arg.ToAccountId {
			result.FromAccount, result.ToAccount, err = AddMoney(ctx, q, txName, arg.FromAccountId, -arg.Amount, arg.ToAccountId, arg.Amount)
		} else {
			result.FromAccount, result.ToAccount, err = AddMoney(ctx, q, txName, arg.ToAccountId, arg.Amount, arg.FromAccountId, -arg.Amount)
		}

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func AddMoney(
	ctx context.Context,
	q *Queries,
	txName any,
	accoun1Id,
	amount1,
	account2Id,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	fmt.Println(txName, "add balance account 1")
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:      accoun1Id,
		Balance: amount1,
	})
	if err != nil {
		return
	}

	fmt.Println(txName, "add balance account 2")
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:      account2Id,
		Balance: amount2,
	})
	if err != nil {
		return
	}

	return
}
