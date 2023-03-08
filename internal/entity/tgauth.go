package entity

type TgAuth struct {
	ID         int   `db:"id"`
	UserID     int   `db:"user_id"`
	TelegramID int64 `db:"telegram_id"`
}
