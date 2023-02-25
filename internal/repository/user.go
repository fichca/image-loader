package repository

import (
	"context"
	"fmt"
	"github.com/fichca/image-loader/internal/config"
	"github.com/fichca/image-loader/internal/entity"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db  *sqlx.DB
	cfg *config.DB
}

func NewUserRepo(db *sqlx.DB, cfg *config.DB) *UserRepo {
	return &UserRepo{
		db:  db,
		cfg: cfg,
	}
}

func (u *UserRepo) Add(ctx context.Context, user entity.User) error {
	query := `INSERT INTO users(name, description, login, password) VALUES (:name, :description, :login, :password)`

	_, err := u.db.NamedExecContext(ctx, query, &user)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

func (u *UserRepo) GetById(ctx context.Context, id int) (entity.User, error) {
	query := `SELECT * FROM users WHERE id = $1`

	var us entity.User

	row := u.db.QueryRowxContext(ctx, query, id)

	err := row.StructScan(&us)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to scan struct user: %w", err)
	}

	return us, nil
}

func (u *UserRepo) Update(ctx context.Context, user entity.User) error {
	query := `UPDATE users set (name, description, login, password) = (:name, :description, :login, :password) WHERE id = :id`

	_, err := u.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

func (u *UserRepo) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepo) GetAll(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepo) RunMigrations() error {
	driver, err := postgres.WithInstance(u.db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to get migration tool driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		u.cfg.Driver, driver)
	if err != nil {
		return fmt.Errorf("failed to connect migration tool: %w", err)
	}

	err = m.Up()
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
