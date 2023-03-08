package repository

import (
	"context"
	"fmt"
	"github.com/fichca/image-loader/internal/entity"
	"github.com/jmoiron/sqlx"
)

type TgAuthRepo struct {
	db *sqlx.DB
}

func NewTgAuthRepo(db *sqlx.DB) *TgAuthRepo {
	return &TgAuthRepo{
		db: db,
	}
}

func (t TgAuthRepo) Authorize(ctx context.Context, userID int, telegramID int64) error {
	query := `INSERT INTO tg_auth(user_id, telegram_id) VALUES ($1, $2)`

	_, err := t.db.ExecContext(ctx, query, userID, telegramID)
	if err != nil {
		return fmt.Errorf("failed to tg auth: %w", err)
	}

	return nil
}

func (t TgAuthRepo) CheckTgAuth(ctx context.Context, tgID int64) (userID int, err error) {
	query := `SELECT * 
              FROM tg_auth tga 
              WHERE tga.telegram_id = $1`

	var tgAuth entity.TgAuth

	row := t.db.QueryRowxContext(ctx, query, tgID)

	err = row.StructScan(&tgAuth)
	if err != nil {
		return 0, fmt.Errorf("failed to tg auth: %w", err)
	}

	return tgAuth.UserID, nil
}
