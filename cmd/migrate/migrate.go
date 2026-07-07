package migrate

import (
	"database/sql"
	"fmt"

	"github.com/goggle-source/authLotServic/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(cfg *config.Cfg) error {

	conn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Db.User, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port, cfg.Db.DbName)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	if errSource, errDB := m.Close(); errSource != nil || errDB != nil {
		_ = db.Close()
		return fmt.Errorf("%w, %w", errDB, errSource)
	}

	if err = db.Close(); err != nil {
		return err
	}

	return nil
}
