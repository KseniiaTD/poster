package postgres

import (
	"context"
	"database/sql"

	"github.com/KseniiaTD/poster/graph/model"
)

func (db *postgresDB) CreateUser(ctx context.Context, newUser model.NewUserInput) (int, error) {
	var row sql.NullInt32

	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	insertUserStmt := `INSERT INTO public.users(login, name, surname, phone, email) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err = tx.QueryRow(insertUserStmt, newUser.Login, newUser.Name, newUser.Surname, newUser.Phone, newUser.Email).Scan(&row)
	if err != nil {
		return 0, err
	}

	userId := row.Int32

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return int(userId), nil

}
