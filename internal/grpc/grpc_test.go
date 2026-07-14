package grpc

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/go-playground/validator/v10"
	auth "github.com/goggle-source/authLotProto/gen/go/auth"
	"github.com/goggle-source/authLotServic/domain"
	mock_servic "github.com/goggle-source/authLotServic/internal/mock"
	"github.com/goggle-source/authLotServic/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestServerAPI_Register(t *testing.T) {
	testName := "Test User"
	testEmail := "test@example.com"
	testPassword := "secret123"
	testToken := "jwt-token"

	tests := []struct {
		name          string
		req           *auth.RegisterUserRequest
		setupMocks    func(authSvc *mock_servic.MockAuthService)
		expectedToken string
		expectedErr   error
		expectedCode  codes.Code
		containsErr   string
		hasError      bool
	}{
		{
			name: "successful_register",
			req: &auth.RegisterUserRequest{
				Name:     testName,
				Email:    testEmail,
				Password: testPassword,
			},
			setupMocks: func(authSvc *mock_servic.MockAuthService) {
				authSvc.On("Register", mock.Anything, models.UserRegister{
					Name:     testName,
					Email:    testEmail,
					Password: testPassword,
				}).Return(testToken, nil)
			},
			expectedToken: testToken,
			expectedErr:   nil,
			expectedCode:  codes.OK,
			hasError:      false,
		},
		{
			name: "validation_error_min",
			req: &auth.RegisterUserRequest{
				Name:     "jonn",
				Email:    "example@gmail.com",
				Password: "passsword_123",
			},
			setupMocks:    nil,
			expectedToken: "",
			expectedErr:   nil,
			expectedCode:  codes.InvalidArgument,
			containsErr:   "error min:",
			hasError:      true,
		},
		{
			name: "validation_error_max",
			req: &auth.RegisterUserRequest{
				Name:     "jonnaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				Email:    "example@gmail.com",
				Password: "passsword_123",
			},
			setupMocks:    nil,
			expectedToken: "",
			expectedErr:   nil,
			expectedCode:  codes.InvalidArgument,
			containsErr:   "error max:",
			hasError:      true,
		},
		{
			name: "validation_error_email",
			req: &auth.RegisterUserRequest{
				Name:     "jonnsina",
				Email:    "exgmail.com",
				Password: "passsword_123",
			},
			setupMocks:    nil,
			expectedToken: "",
			expectedErr:   nil,
			expectedCode:  codes.InvalidArgument,
			containsErr:   "email is invalid:",
			hasError:      true,
		},
		{
			name: "auth_service_error",
			req: &auth.RegisterUserRequest{
				Name:     testName,
				Email:    testEmail,
				Password: testPassword,
			},
			setupMocks: func(authSvc *mock_servic.MockAuthService) {
				authSvc.On("Register", mock.Anything, mock.Anything).Return("", domain.ErrGenerateJWT)
			},
			expectedToken: "",
			expectedErr:   domain.ErrGenerateJWT,
			expectedCode:  codes.Internal,
			hasError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuth := new(mock_servic.MockAuthService)

			if tt.setupMocks != nil {
				tt.setupMocks(mockAuth)
			}

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			validate := validator.New()
			srv := &ServerAPI{
				log:      logger,
				auth:     mockAuth,
				validate: validate,
			}

			resp, err := srv.Register(context.Background(), tt.req)

			if tt.hasError {
				require.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedCode, st.Code())

				if tt.containsErr != "" {
					require.Contains(t, err.Error(), tt.containsErr)
				}
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, tt.expectedToken, resp.Token)
			}

			mockAuth.AssertExpectations(t)
		})
	}
}

func TestServerAPI_Login(t *testing.T) {
	testEmail := "login@example.com"
	testPassword := "validPass123"
	testName := "Alice"
	testToken := "jwt-token-xyz"

	tests := []struct {
		name          string
		req           *auth.LoginUserRequest
		setupMocks    func(authSvc *mock_servic.MockAuthService)
		expectedName  string
		expectedToken string
		expectedCode  codes.Code
		hasError      bool
		containsErr   string
	}{
		{
			name: "successful_login",
			req: &auth.LoginUserRequest{
				Email:    testEmail,
				Password: testPassword,
			},
			setupMocks: func(authSvc *mock_servic.MockAuthService) {
				authSvc.On("Login", mock.Anything, models.UserLogin{
					Email:    testEmail,
					Password: testPassword,
				}).Return(testName, testToken, nil)
			},
			expectedName:  testName,
			expectedToken: testToken,
			hasError:      false,
			expectedCode:  codes.OK,
		},
		{
			name: "validation_error_min",
			req: &auth.LoginUserRequest{
				Email:    "example@gmail.com",
				Password: "pa",
			},
			setupMocks:    nil,
			expectedToken: "",
			expectedCode:  codes.InvalidArgument,
			containsErr:   "error min:",
			hasError:      true,
		},
		{
			name: "validation_error_max",
			req: &auth.LoginUserRequest{
				Email:    "example@gmail.com",
				Password: "passsword_123asdddddddddddddddddddddddddddddddddddddddddddddddddddddddddssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssddddssssssssssssssssssssssssssssssssdafdwefgdsgzdfhbdzhjdzghnfsGddsdccssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssss",
			},
			setupMocks:    nil,
			expectedToken: "",
			expectedCode:  codes.InvalidArgument,
			containsErr:   "error max:",
			hasError:      true,
		},
		{
			name: "validation_error_email",
			req: &auth.LoginUserRequest{
				Email:    "exgmail.com",
				Password: "passsword_123",
			},
			setupMocks:    nil,
			expectedToken: "",
			expectedCode:  codes.InvalidArgument,
			containsErr:   "email is invalid:",
			hasError:      true,
		},
		{
			name: "validation_error",
			req: &auth.LoginUserRequest{
				Email:    "invalid",
				Password: "",
			},
			setupMocks: func(authSvc *mock_servic.MockAuthService) {
			},
			expectedName:  "",
			expectedToken: "",
			expectedCode:  codes.InvalidArgument,
			hasError:      true,
		},
		{
			name: "auth_login_error",
			req: &auth.LoginUserRequest{
				Email:    testEmail,
				Password: testPassword,
			},
			setupMocks: func(authSvc *mock_servic.MockAuthService) {
				authSvc.On("Login", mock.Anything, mock.Anything).Return("", "", domain.ErrGenerateJWT)
			},
			expectedName:  "",
			expectedToken: "",
			expectedCode:  codes.Internal,
			hasError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuth := new(mock_servic.MockAuthService)

			if tt.setupMocks != nil {
				tt.setupMocks(mockAuth)
			}

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			validate := validator.New()
			srv := &ServerAPI{
				log:      logger,
				validate: validate,
				auth:     mockAuth,
			}

			resp, err := srv.Login(context.Background(), tt.req)

			if tt.hasError {
				require.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedCode, st.Code())
				if tt.containsErr != "" {
					require.Contains(t, err.Error(), tt.containsErr)
				}
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, tt.expectedToken, resp.Token)
				require.Equal(t, tt.expectedName, resp.Name)
			}

			mockAuth.AssertExpectations(t)
		})
	}
}
