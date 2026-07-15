package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/goggle-source/authLotServic/domain"
	"github.com/goggle-source/authLotServic/internal/config"
	"github.com/goggle-source/authLotServic/internal/metric"
	"github.com/goggle-source/authLotServic/internal/models"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Db struct {
	DB *sql.DB
}

func Init(cfg *config.Cfg) *Db {
	conn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Db.User, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port, cfg.Db.DbName)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(cfg.Db.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Db.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.Db.ConnMaxIdleTime)
	db.SetConnMaxLifetime(cfg.Db.ConnMaxLifeTime)

	return &Db{
		DB: db,
	}
}

func (d *Db) Register(ctx context.Context, userAddDatabase models.UserAddDatabase) error {
	const op = "repository.Register"

	_, err := d.DB.ExecContext(ctx, `INSERT INTO users (userName, email, pass_hash, uid) VALUES 
	($1, $2, $3, $4)`, userAddDatabase.Name, userAddDatabase.Email,
		userAddDatabase.PasswordHash, userAddDatabase.Id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return fmt.Errorf("%s:%w", op, domain.ErrEmail)
			}
		}
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (d *Db) Login(ctx context.Context, userValidateInDatabase models.UserValidateInDatabase) (string, string, []byte, error) {
	const op = "repository.Login"

	var name, id string
	var passHash []byte
	err := d.DB.QueryRowContext(ctx, "SELECT userName, uid,  pass_hash FROM users WHERE email = $1", userValidateInDatabase.Email).Scan(&name, &id, &passHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", []byte{}, fmt.Errorf("%s:%w", op, domain.ErrUserNoFound)
		}
		return "", "", []byte{}, fmt.Errorf("%s:%w", op, err)
	}

	return name, id, passHash, nil
}

func (d *Db) HealthCheack(ctx context.Context) (metric.DBMetric, error) {
	const op = "repository.Check"

	var result metric.DBMetric
	if err := d.DB.Ping(); err != nil {
		result.ConnDB = false
		return result, fmt.Errorf("%s:%w", op, err)
	}

	_, err := d.DB.ExecContext(ctx, "SELECT 1")
	if err != nil {
		result.ConnDB = false
		return result, fmt.Errorf("%s:%w", op, err)
	}

	err = d.DB.QueryRowContext(ctx, "SELECT count(*) FROM pg_stat_activity WHERE state = `active` ").Scan(&result.ActiveConnection)
	if err != nil {
		result.ActiveConnection = 0
		return result, fmt.Errorf("%s:%w", op, err)
	}

	err = d.DB.QueryRowContext(ctx, "SELECT count(*) FROM pg_stat_activity").Scan(&result.CountConnection)
	if err != nil {
		result.CountConnection = 0
		return result, fmt.Errorf("%s:%w", op, err)
	}

	err = d.DB.QueryRowContext(ctx, "SELECT * FROM pg_stat_wal()").Scan(&result.CountMemory)
	if err != nil {
		result.CountMemory = 0
		return result, fmt.Errorf("%s:%w", op, err)
	}

	return result, nil
}
