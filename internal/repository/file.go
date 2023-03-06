package repository

import (
	"context"
	"fmt"
	"github.com/fichca/image-loader/internal/entity"
	"github.com/jmoiron/sqlx"
)

type ImageRepo struct {
	db *sqlx.DB
}

func NewImageRepo(db *sqlx.DB) *ImageRepo {
	return &ImageRepo{
		db: db,
	}
}

func (i *ImageRepo) Add(ctx context.Context, image entity.Image) error {
	query := `INSERT INTO images(user_id, name, extension) VALUES (:user_id, :name, :extension)`

	_, err := i.db.NamedExecContext(ctx, query, &image)
	if err != nil {
		return fmt.Errorf("failed to insert image: %w", err)
	}

	return nil
}

func (i *ImageRepo) GetById(ctx context.Context, id int) (entity.Image, error) {
	query := `SELECT * FROM images WHERE id = $1`

	var img entity.Image

	row := i.db.QueryRowxContext(ctx, query, id)

	err := row.StructScan(&img)
	if err != nil {
		return entity.Image{}, fmt.Errorf("failed to scan struct image: %w", err)
	}

	return img, nil
}

func (i *ImageRepo) GetAllByUserId(ctx context.Context, userID int) ([]entity.Image, error) {
	query := `SELECT * FROM images WHERE user_id = $1`

	var images []entity.Image

	rows, err := i.db.QueryxContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query images: %w", err)
	}

	for rows.Next() {
		var img entity.Image

		err := rows.StructScan(&img)
		if err != nil {
			return nil, err
		}

		images = append(images, img)
	}

	return images, nil
}
