package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/dxtym/bankrupt/utils"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := utils.HashPassword(utils.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       utils.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

// test create user
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

// test get user
func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	user := createRandomUser(t)
	newFullName := utils.RandomOwner()
	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: user.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.Equal(t, user.Username, updatedUser.Username)
	require.Equal(t, user.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newFullName, updatedUser.FullName)
	require.Equal(t, user.Email, updatedUser.Email)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	user := createRandomUser(t)
	newEmail := utils.RandomEmail()
	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: user.Username,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.Equal(t, user.Username, updatedUser.Username)
	require.Equal(t, user.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, user.FullName, updatedUser.FullName)
	require.Equal(t, newEmail, updatedUser.Email)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	user := createRandomUser(t)
	newPassword, _ := utils.HashPassword(utils.RandomString(6))
	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: user.Username,
		HashedPassword: sql.NullString{
			String: newPassword,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.Equal(t, user.Username, updatedUser.Username)
	require.Equal(t, newPassword, updatedUser.HashedPassword)
	require.Equal(t, user.FullName, updatedUser.FullName)
	require.Equal(t, user.Email, updatedUser.Email)
}

func TestUpdateUserAllFields(t *testing.T) {
	user := createRandomUser(t)
	newFullName := utils.RandomOwner()
	newEmail := utils.RandomEmail()
	newPassword, _ := utils.HashPassword(utils.RandomString(6))
	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: user.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		HashedPassword: sql.NullString{
			String: newPassword,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.Equal(t, user.Username, updatedUser.Username)
	require.Equal(t, newPassword, updatedUser.HashedPassword)
	require.NotEqual(t, user.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newFullName, updatedUser.FullName)
	require.NotEqual(t, user.FullName, updatedUser.FullName)
	require.Equal(t, newEmail, updatedUser.Email)
	require.NotEqual(t, user.Email, updatedUser.Email)
}
