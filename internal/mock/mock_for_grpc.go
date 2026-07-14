package mock

import (
	"context"

	"github.com/goggle-source/authLotServic/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, user models.UserRegister) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, userLogin models.UserLogin) (name string, token string, err error) {
	args := m.Called(ctx, userLogin)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockAuthService) HealthyCheack(ctx context.Context) (map[string]string, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockAuthService) ValidateUser(ctx context.Context, userID string) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}
