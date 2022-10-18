package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"guthub.com/Edbeer/microservices/internal/core"
)

func Test_Create(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	userStorage := newAccountStorage(sqlxDB)

	t.Run("Create", func(t *testing.T) {
		columns := []string{
			"name",
			"email",
			"password",
			"role",
		}
		rows := sqlmock.NewRows(columns).AddRow(
			"PavelV",
			"edbeer123mtn@gmail.com",
			"12345678",
			"user",
		)

		user := &core.User{
			Name:  "PavelV",
			Email: "edbeer123mtn@gmail.com",
			Pass:  "12345678",
			Role:  "user",
		}

		query := `INSERT INTO users (name, email, password, role, created_at) 
			VALUES ($1, $2, $3, $4, now()) 
			RETURNING *`
		mock.ExpectQuery(query).WithArgs(
			&user.Name, &user.Email, &user.Pass, &user.Role,
		).WillReturnRows(rows)

		createdUser, err := userStorage.Create(context.Background(), user)
		require.NoError(t, err)
		require.NotNil(t, createdUser)
		require.Equal(t, createdUser, user)
	})
}

func Test_FindUserByEmail(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	userStorage := newAccountStorage(sqlxDB)

	t.Run("FindUserByEmail", func(t *testing.T) {
		uid := uuid.New()

		columns := []string{
			"user_id",
			"name",
			"email",
		}
		rows := sqlmock.NewRows(columns).AddRow(
			uid,
			"PavelV",
			"edbeer123mtn@gmail.com",
		)

		testUser := &core.User{
			Uuid:    uid,
			Name:  "PavelV",
			Email: "edbeer123mtn@gmail.com",
		}

		query := `SELECT user_id, name, email, password, role, created_at
			FROM users
			WHERE email = $1`
		mock.ExpectQuery(query).WithArgs(&testUser.Email).WillReturnRows(rows)

		user, err := userStorage.FindByEmail(context.Background(), testUser)
		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, user.Name, testUser.Name)
	})
}

func Test_GetUserByID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	userStorage := newAccountStorage(sqlxDB)

	t.Run("GetUserByID", func(t *testing.T) {
		uid := uuid.New()

		columns := []string{
			"user_id",
			"name",
			"email",
		}
		rows := sqlmock.NewRows(columns).AddRow(
			uid,
			"PavelV",
			"edbeer123mtn@gmail.com",
		)

		testUser := &core.User{
			Uuid:  uid,
			Name:  "PavelV",
			Email: "edbeer123mtn@gmail.com",
		}

		query := `SELECT user_id, name, email, password, role, created_at
			FROM users
			WHERE user_id = $1`
		mock.ExpectQuery(query).WithArgs(uid).WillReturnRows(rows)

		user, err := userStorage.GetUserByID(context.Background(), uid)
		require.NoError(t, err)
		require.Equal(t, user.Name, testUser.Name)
		fmt.Printf("test user: %s \n", testUser.Name)
		fmt.Printf("user: %s \n", user.Name)
	})
}
