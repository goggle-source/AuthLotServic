package servic

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/goggle-source/authLotServic/domain"
	"github.com/goggle-source/authLotServic/internal/config"
	mock_servic "github.com/goggle-source/authLotServic/internal/mock"
	"github.com/goggle-source/authLotServic/internal/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestServicApp_Register(t *testing.T) {

	testEmail := "test@example.com"
	testName := "John Doe"
	testPassword := "secret123"

	tests := []struct {
		name          string
		userRegister  models.UserRegister
		setupMocks    func(mockDB *mock_servic.MockDatabase, mockTG *mock_servic.MockTokenGenerator)
		expectedToken string
		expectedErr   error
	}{
		{
			name: "successfull call",
			userRegister: models.UserRegister{
				Email:    testEmail,
				Password: testPassword,
				Name:     testName,
			},
			setupMocks: func(mockDB *mock_servic.MockDatabase, mockTG *mock_servic.MockTokenGenerator) {
				mockDB.On("Register", mock.Anything, mock.MatchedBy(func(user models.UserAddDatabase) bool {
					return user.Email == testEmail &&
						user.Name == testName &&
						len(user.PasswordHash) > 0 &&
						user.Id != ""
				})).Return(nil)

				mockTG.On("GenerateJWTToken", mock.Anything, mock.AnythingOfType("string")).
					Return("jwt-token-123", nil)
			},
			expectedToken: "jwt-token-123",
			expectedErr:   nil,
		},
		{
			name: "error in database",
			userRegister: models.UserRegister{
				Email:    testEmail,
				Password: testPassword,
				Name:     testName,
			},
			setupMocks: func(mockDB *mock_servic.MockDatabase, mockTG *mock_servic.MockTokenGenerator) {
				mockDB.On("Register", mock.Anything, mock.Anything).Return(domain.ErrEmail)
			},
			expectedToken: "",
			expectedErr:   domain.ErrEmail,
		},
		{
			name: "err generate JWT",
			userRegister: models.UserRegister{
				Email:    testEmail,
				Password: testPassword,
				Name:     testName,
			},
			setupMocks: func(mockDB *mock_servic.MockDatabase, mockTG *mock_servic.MockTokenGenerator) {
				mockDB.On("Register", mock.Anything, mock.Anything).Return(nil)
				mockTG.On("GenerateJWTToken", mock.Anything, mock.Anything).
					Return("", domain.ErrGenerateJWT)
			},
			expectedToken: "",
			expectedErr:   domain.ErrGenerateJWT,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB := new(mock_servic.MockDatabase)
			mockTG := new(mock_servic.MockTokenGenerator)

			if tt.setupMocks != nil {
				tt.setupMocks(mockDB, mockTG)
			}
			log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

			svc := Init(log, mockDB, mockTG)

			gotToken, gotErr := svc.Register(context.Background(), tt.userRegister)

			if tt.expectedErr != nil {
				require.Error(t, gotErr)
				require.Contains(t, gotErr.Error(), tt.expectedErr.Error())
			} else {
				require.NoError(t, gotErr)
				require.Equal(t, tt.expectedToken, gotToken)
			}

			mockDB.AssertExpectations(t)
			mockTG.AssertExpectations(t)
		})
	}
}

func TestServicApp_Login(t *testing.T) {
	testEmail := "login@example.com"
	testPassword := "validPass123"
	testName := "Alice"
	testID := "user-id-123"
	testToken := "jwt-token-xyz"

	validHash, _ := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)

	tests := []struct {
		name          string
		userLogin     models.UserLogin
		setupMocks    func(mockDB *mock_servic.MockDatabase, mockTG *mock_servic.MockTokenGenerator)
		expectedName  string
		expectedToken string
		expectedErr   error
	}{
		{
			name: "successful_login",
			userLogin: models.UserLogin{
				Email:    testEmail,
				Password: testPassword,
			},
			setupMocks: func(mockDB *mock_servic.MockDatabase, mockTG *mock_servic.MockTokenGenerator) {
				mockDB.On("Login", mock.Anything, models.UserValidateInDatabase{Email: testEmail}).
					Return(testName, testID, validHash, nil)
				mockTG.On("GenerateJWTToken", mock.Anything, testID).
					Return(testToken, nil)
			},
			expectedName:  testName,
			expectedToken: testToken,
			expectedErr:   nil,
		},
		{
			name: "database_login_error",
			userLogin: models.UserLogin{
				Email:    testEmail,
				Password: testPassword,
			},
			setupMocks: func(mockDB *mock_servic.MockDatabase, mockTG *mock_servic.MockTokenGenerator) {
				mockDB.On("Login", mock.Anything, models.UserValidateInDatabase{Email: testEmail}).
					Return("", "", []byte(nil), domain.ErrUserNoFound)
			},
			expectedName:  "",
			expectedToken: "",
			expectedErr:   domain.ErrUserNoFound,
		},
		{
			name: "invalid_password",
			userLogin: models.UserLogin{
				Email:    testEmail,
				Password: "wrongpass",
			},
			setupMocks: func(mockDB *mock_servic.MockDatabase, mockTG *mock_servic.MockTokenGenerator) {
				mockDB.On("Login", mock.Anything, models.UserValidateInDatabase{Email: testEmail}).
					Return(testName, testID, validHash, nil)
			},
			expectedName:  "",
			expectedToken: "",
			expectedErr:   domain.ErrPasswordOrEmail,
		},
		{
			name: "token_generation_error",
			userLogin: models.UserLogin{
				Email:    testEmail,
				Password: testPassword,
			},
			setupMocks: func(mockDB *mock_servic.MockDatabase, mockTG *mock_servic.MockTokenGenerator) {
				mockDB.On("Login", mock.Anything, models.UserValidateInDatabase{Email: testEmail}).
					Return(testName, testID, validHash, nil)
				mockTG.On("GenerateJWTToken", mock.Anything, testID).
					Return("", domain.ErrGenerateJWT)
			},
			expectedName:  "",
			expectedToken: "",
			expectedErr:   domain.ErrGenerateJWT,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(mock_servic.MockDatabase)
			mockTG := new(mock_servic.MockTokenGenerator)

			if tt.setupMocks != nil {
				tt.setupMocks(mockDB, mockTG)
			}
			log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

			svc := Init(log, mockDB, mockTG)

			gotName, gotToken, gotErr := svc.Login(context.Background(), tt.userLogin)

			if tt.expectedErr != nil {
				require.Error(t, gotErr)
				require.ErrorIs(t, gotErr, tt.expectedErr)
			} else {
				require.NoError(t, gotErr)
				require.Equal(t, tt.expectedName, gotName)
				require.Equal(t, tt.expectedToken, gotToken)
			}

			mockDB.AssertExpectations(t)
			mockTG.AssertExpectations(t)
		})
	}
}

func GetPathKey() string {
	cfg := config.Load()

	return cfg.Path
}
