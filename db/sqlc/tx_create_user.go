package db

import (
	"context"
)

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error // callback function to run after creating the user
}

type CreateUserTxResult struct {
	User User
}

func (store *SqlStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var txResult CreateUserTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txResult.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(txResult.User)
	})

	return txResult, err
}
