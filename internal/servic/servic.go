package servic

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/goggle-source/authLotServic/domain"
	"github.com/goggle-source/authLotServic/internal/metric"
	"github.com/goggle-source/authLotServic/internal/models"
	"github.com/google/uuid"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"golang.org/x/crypto/bcrypt"
)

const (
	cost = 15
)

type Database interface {
	Register(ctx context.Context, userAddDatabase models.UserAddDatabase) error
	Login(ctx context.Context, userValidateInDatabase models.UserValidateInDatabase) (name string, uid string, err error)
	HealthCheack(ctx context.Context) (metric.DBMetric, error)
	ValidateUserId(ctx context.Context, id string) (bool, error)
}

type ServicApp struct {
	d        Database
	tokenSSL *rsa.PrivateKey
}

func Init(log *slog.Logger, d Database, tokenSSL *rsa.PrivateKey) *ServicApp {
	return &ServicApp{
		d:        d,
		tokenSSL: tokenSSL,
	}
}

func (s *ServicApp) Register(ctx context.Context, userRegister models.UserRegister) (token string, err error) {
	const op = "servic.Register"

	uid := uuid.New()
	id := uid.String()

	bytes, err := bcrypt.GenerateFromPassword([]byte(userRegister.Password), cost)
	if err != nil {
		return "", fmt.Errorf("%s error generate hash for password:%w", op, err)
	}

	userAddDatabase := models.UserAddDatabase{
		Email:        userRegister.Email,
		PasswordHash: bytes,
		Name:         userRegister.Name,
		Id:           uid.String(),
	}

	err = s.d.Register(ctx, userAddDatabase)
	if err != nil {
		return "", fmt.Errorf("%s err register in database layer:%w", op, err)
	}

	token, err = GenerateJWTToken(ctx, id, s.tokenSSL)
	if err != nil {
		return "", fmt.Errorf("%s:%w", err, domain.ErrGenerateJWT)
	}

	return token, nil
}

func (s *ServicApp) Login(ctx context.Context, userLogin models.UserLogin) (name string, token string, err error) {
	const op = "servic.Login"

	userValidInDatabase := models.UserValidateInDatabase{
		Email:    userLogin.Email,
		Password: userLogin.Password,
	}

	name, id, err := s.d.Login(ctx, userValidInDatabase)
	if err != nil {
		return "", "", fmt.Errorf("%s err login in database layer:%w", op, err)
	}

	token, err = GenerateJWTToken(ctx, id, s.tokenSSL)
	if err != nil {
		return "", "", fmt.Errorf("%s:%w", err, domain.ErrGenerateJWT)
	}

	return name, token, nil
}

func (s *ServicApp) HealthyCheack(ctx context.Context) (map[string]string, error) {

	details := make(map[string]string)

	detailsDB, err := s.d.HealthCheack(ctx)

	counts, err := cpu.Percent(1*time.Second, false)
	if err != nil {
		details["errCpu"] = "error getting CPU load"
	} else {
		details["cpu"] = strconv.Itoa(int(counts[0]))
	}
	memory, err := mem.VirtualMemory()
	if err != nil {
		details["errMemory"] = "error getting memory usage"
	} else {
		details["Totalmemory"] = strconv.Itoa(int(memory.Total))
		details["UseMemory"] = strconv.Itoa(int(memory.Used))
		details["PercentMemory"] = strconv.Itoa(int(memory.UsedPercent))
	}

	if detailsDB.ConnDB {
		details["pingDB"] = "true"
	} else {
		details["pingDB"] = "false"
	}
	details["ActiveConnDB"] = strconv.Itoa(detailsDB.ActiveConnection)
	details["CountConnDB"] = strconv.Itoa(detailsDB.CountConnection)
	details["CountMemoryDB"] = strconv.Itoa(detailsDB.CountMemory)

	return details, nil
}

func (s *ServicApp) ValidateUser(ctx context.Context, id string) (bool, error) {
	const op = "servic.ValidateUser"

	isValid, err := s.d.ValidateUserId(ctx, id)
	if err != nil {
		return false, fmt.Errorf("%s err validateUserID in database layer:%w", op, err)
	}

	return isValid, err
}
