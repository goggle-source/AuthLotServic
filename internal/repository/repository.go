package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/goggle-source/authLotServic/internal/config"
	"github.com/goggle-source/authLotServic/internal/lib/logger"
	"github.com/goggle-source/authLotServic/internal/metric"
	"github.com/goggle-source/authLotServic/internal/models"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Db struct {
	DB  *sql.DB
	log *slog.Logger
}

func Init(cfg *config.Cfg, log *slog.Logger) *Db {
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
		DB:  db,
		log: log,
	}
}

func (d *Db) Register(ctx context.Context, userAddDatabase models.UserAddDatabase) error {
	const op = "repository.Register"

	log := d.log.With(slog.String("op", op))

	log.Info("start register user")

	rows := d.DB.QueryRowContext(ctx, `INSERT INTO users (userName, email, pass_hash, uid) VALUES 
	($1, $2, $3, $4)`, userAddDatabase.Name, userAddDatabase.Email,
		userAddDatabase.PasswordHash, userAddDatabase.Id)
	if rows.Err() != nil {
		log.Error("error add user in database", logger.Err(rows.Err()))
		return ValidateErrorsPostgresql(rows.Err())
	}

	log.Info("success register user")

	return nil
}

func (d *Db) Login(ctx context.Context, userValidateInDatabase models.UserValidateInDatabase) (string, string, error) {
	const op = "repository.Login"

	log := d.log.With(slog.String("op", op))

	log.Info("start login user")

	var name, id string
	var passHash []byte
	err := d.DB.QueryRowContext(ctx, "SELECT userName, uid,  pass_hash FROM users WHERE email = $1", userValidateInDatabase.Email).Scan(&name, &id, &passHash)
	if err != nil {
		log.Error("error get user for database", logger.Err(err))
		return "", "", ValidateErrorsPostgresql(err)
	}

	if err := bcrypt.CompareHashAndPassword(passHash, []byte(userValidateInDatabase.Password)); err != nil {
		log.Error("password not equal to password from database", logger.Err(err))
		return "", "", ErrPassword
	}
	log.Info("success Login user")

	return name, id, nil
}

func (d *Db) HealthCheack(ctx context.Context) (metric.DBMetric, error) {
	const op = "repository.Check"

	log := d.log.With(slog.String("op", op))

	log.Info("start healthy")

	var result metric.DBMetric
	if err := d.DB.Ping(); err != nil {
		log.Error("the database is not responding", logger.Err(err))
		result.ConnDB = false
		return result, ValidateErrorsPostgresql(err)
	}

	_, err := d.DB.ExecContext(ctx, "SELECT 1")
	if err != nil {
		log.Error("the database did not complete the request", logger.Err(err))
		result.ConnDB = false
		return result, ValidateErrorsPostgresql(err)
	}

	err = d.DB.QueryRowContext(ctx, "SELECT count(*) FROM pg_stat_activity WHERE state = `active`").Scan(result.ActiveConnection)
	if err != nil {
		log.Error("couldn't get the number of active connections", logger.Err(err))
		result.ActiveConnection = 0
		return result, ValidateErrorsPostgresql(err)
	}

	err = d.DB.QueryRowContext(ctx, "SELECT count(*) FROM pg_stat_activity").Scan(&result.CountConnection)
	if err != nil {
		log.Error("couldn't get the number of all connections", logger.Err(err))
		result.CountConnection = 0
		return result, ValidateErrorsPostgresql(err)
	}

	err = d.DB.QueryRowContext(ctx, "SELECT * FROM pg_stat_wal()").Scan(&result.CountMemory)
	if err != nil {
		log.Error("couldn't get the number memory", logger.Err(err))
		result.CountMemory = 0
		return result, ValidateErrorsPostgresql(err)
	}
	log.Info("success healthycheack")

	return result, nil
}

func (d *Db) ValidateUserId(ctx context.Context, id int) (bool, error) {
	const op = "repository.ValidateUserId"

	log := d.log.With(slog.String("op", op))

	log.Info("start validateUserID")

	var email string

	err := d.DB.QueryRowContext(ctx, "SELECT email FROM users WHERE uid = $1", id).Scan(&email)
	if err != nil {
		log.Error("error get email users", logger.Err(err))
		return false, ValidateErrorsPostgresql(err)
	}

	log.Info("success validate userID")

	if email != "" {
		return true, nil
	} else {
		return false, nil
	}
}
