package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/goggle-source/authLotServic/domain"
	"github.com/goggle-source/authLotServic/internal/models"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestDb_Register(t *testing.T) {
	type args struct {
		ctx  context.Context
		user models.UserAddDatabase
	}
	tests := []struct {
		name    string
		args    args
		mock    func(mock sqlmock.Sqlmock)
		wantErr error
	}{
		{
			name: "successful registration",
			args: args{
				ctx: context.Background(),
				user: models.UserAddDatabase{
					Name:         "John",
					Email:        "john@example.com",
					PasswordHash: []byte("pass"),
					Id:           "123",
				},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO users (userName, email, pass_hash, uid) VALUES ($1, $2, $3, $4)`).
					WithArgs("John", "john@example.com", []byte("pass"), "123").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: nil,
		},
		{
			name: "duplicate email error",
			args: args{
				ctx: context.Background(),
				user: models.UserAddDatabase{
					Name:         "Jane",
					Email:        "jane@example.com",
					PasswordHash: []byte("pass"),
					Id:           "456",
				},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO users (userName, email, pass_hash, uid) VALUES ($1, $2, $3, $4)`).
					WithArgs("Jane", "jane@example.com", []byte("pass"), "456").
					WillReturnError(&pq.Error{Code: "23505"})
			},
			wantErr: domain.ErrEmail,
		},
		{
			name: "other database error",
			args: args{
				ctx: context.Background(),
				user: models.UserAddDatabase{
					Name:         "Bob",
					Email:        "bob@example.com",
					PasswordHash: []byte("pass"),
					Id:           "789",
				},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO users (userName, email, pass_hash, uid) VALUES ($1, $2, $3, $4)`).
					WithArgs("Bob", "bob@example.com", []byte("pass"), "789").
					WillReturnError(errors.New("connection refused"))
			},
			wantErr: errors.New("connection refused"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			require.NoError(t, err)
			defer db.Close()

			d := &Db{DB: db}
			tt.mock(mock)

			err = d.Register(tt.args.ctx, tt.args.user)

			if tt.wantErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)

				if tt.wantErr == domain.ErrEmail {
					require.True(t, errors.Is(err, domain.ErrEmail), "expected error to be domain.ErrEmail")
				} else {

					require.ErrorContains(t, err, tt.wantErr.Error())
				}
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDb_Login(t *testing.T) {

	type args struct {
		ctx  context.Context
		user models.UserValidateInDatabase
	}
	tests := []struct {
		name     string
		args     args
		mock     func(mock sqlmock.Sqlmock)
		wantName string
		wantID   string
		passHash []byte
		wantErr  error
	}{
		{
			name: "successfull call",
			args: args{
				ctx: context.Background(),
				user: models.UserValidateInDatabase{
					Email: "test@example.com",
				},
			},
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"userName", "uid", "pass_hash"}).
					AddRow("JohnDoe", "123", []byte("password"))
				mock.ExpectQuery("SELECT userName, uid,  pass_hash FROM users WHERE email = $1").
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			wantName: "JohnDoe",
			passHash: []byte("password"),
			wantID:   "123",
			wantErr:  nil,
		},
		{
			name: "пользователь не найден",
			args: args{
				ctx: context.Background(),
				user: models.UserValidateInDatabase{
					Email: "notfound@example.com",
				},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT userName, uid,  pass_hash FROM users WHERE email = $1").
					WithArgs("notfound@example.com").
					WillReturnError(sql.ErrNoRows)
			},
			wantName: "",
			wantID:   "",
			wantErr:  domain.ErrUserNoFound,
		},
		{
			name: "ошибка базы данных",
			args: args{
				ctx: context.Background(),
				user: models.UserValidateInDatabase{
					Email: "error@example.com",
				},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT userName, uid,  pass_hash FROM users WHERE email = $1").
					WithArgs("error@example.com").
					WillReturnError(errors.New("connection lost"))
			},
			wantName: "",
			wantID:   "",
			wantErr:  errors.New("connection lost"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаём мок с точным сравнением SQL (без регулярных выражений)
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			require.NoError(t, err)
			defer db.Close()

			d := &Db{DB: db}
			tt.mock(mock)

			name, id, passHash, err := d.Login(tt.args.ctx, tt.args.user)

			if tt.wantErr == nil {
				require.NoError(t, err)
				require.Equal(t, tt.wantName, name)
				require.Equal(t, tt.wantID, id)
				require.Equal(t, tt.passHash, passHash)
			} else {
				require.Error(t, err)
				switch tt.wantErr {
				case domain.ErrUserNoFound:
					require.True(t, errors.Is(err, domain.ErrUserNoFound))
				case domain.ErrPasswordOrEmail:
					require.True(t, errors.Is(err, domain.ErrPasswordOrEmail))
				default:
					require.ErrorContains(t, err, tt.wantErr.Error())
				}
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
