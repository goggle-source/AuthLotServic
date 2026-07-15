package mock

import (
	"context"

	"github.com/goggle-source/authLotServic/internal/metric"
	"github.com/goggle-source/authLotServic/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) Register(ctx context.Context, user models.UserAddDatabase) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

type MockTokenGenerator struct {
	mock.Mock
}

func (m *MockTokenGenerator) GenerateJWTToken(ctx context.Context, id string) (token string, err error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}

func (m *MockDatabase) Login(ctx context.Context, userValidateInDatabase models.UserValidateInDatabase) (name string, uid string, passHash []byte, err error) {
	args := m.Called(ctx, userValidateInDatabase)
	return args.String(0), args.String(1), args.Get(2).([]byte), args.Error(3)
}

func (m *MockDatabase) HealthCheack(ctx context.Context) (metric.DBMetric, error) {
	args := m.Called(ctx)
	return args.Get(0).(metric.DBMetric), args.Error(1)
}
