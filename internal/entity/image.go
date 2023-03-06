package entity

type Image struct {
	ID        int    `db:"id"`
	UserID    int    `db:"user_id"`
	Name      string `db:"name"`
	Extension string `db:"extension"`
}
