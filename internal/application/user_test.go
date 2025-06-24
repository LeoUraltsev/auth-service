package application

import (
	"context"
	"github.com/LeoUraltsev/auth-service/internal/domain/users"
	mockusers "github.com/LeoUraltsev/auth-service/internal/domain/users/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"log/slog"
	"os"
	"testing"
)

// todo: все тесты, cover >80%
//go test ./internal/application/... -coverprofile=coverage.out
//go tool cover -html=coverage.out

var log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

func TestUserServiceHandler_CreateUser(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := mockusers.NewMockUserRepository(ctrl)
	passwordHasher := mockusers.NewMockPasswordHasher(ctrl)

	type fields struct {
		save         *gomock.Call
		checkEmail   *gomock.Call
		passwordHash *gomock.Call
	}

	type args struct {
		name     string
		email    string
		password string
	}

	cases := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "success creation",
			fields: fields{
				save: repository.EXPECT().
					Save(context.Background(), gomock.Any()).Return(nil).
					AnyTimes(),
				checkEmail: repository.EXPECT().
					ExistsByEmail(context.Background(), gomock.Any()).
					Return(false, nil).
					AnyTimes(),
				passwordHash: passwordHasher.EXPECT().
					Hash(gomock.Any()).
					Return([]byte("hashpassword"), nil).
					AnyTimes(),
			},
			args: args{
				name:     "testname",
				email:    "test@mail.ru",
				password: "testtest",
			},
		},
	}

	for _, tt := range cases {
		service := NewUserService(repository, passwordHasher, log)
		uuid, err := service.CreateUser(context.Background(), tt.args.name, tt.args.email, tt.args.password)

		assert.NoError(t, err, "should not error")
		assert.NotNil(t, uuid, "should return uuid")
		assert.NotEmpty(t, uuid, "uuid should not be empty")

	}

}

func TestUserServiceHandler_checkUniqueEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	repository := mockusers.NewMockUserRepository(ctrl)
	repository.EXPECT().
		ExistsByEmail(context.Background(), gomock.Any()).
		Return(false, nil).
		AnyTimes()
	email, err := users.NewEmail("email@gmail.com")
	assert.NoError(t, err, "should not error")
	service := NewUserService(repository, nil, log)
	err = service.checkUniqueEmail(context.Background(), email)
	assert.NoError(t, err, "should not error")
}

func TestUserServiceHandler_checkUniqueEmail_failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	repository := mockusers.NewMockUserRepository(ctrl)
	repository.EXPECT().
		ExistsByEmail(context.Background(), gomock.Any()).
		Return(true, nil).
		AnyTimes()
	email, err := users.NewEmail("email@gmail.com")
	assert.NoError(t, err, "should not error")
	service := NewUserService(repository, nil, log)
	err = service.checkUniqueEmail(context.Background(), email)
	assert.Error(t, err, "should error")
}

func TestUserServiceHandler_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	email, _ := users.NewEmail("success@email.ru")
	pass, _ := users.NewPassword([]byte("hashpassword"))
	user, err := users.CreateUser("name", email, pass)
	assert.NoError(t, err, "should not error")

	repository := mockusers.NewMockUserRepository(ctrl)
	repository.EXPECT().
		Get(context.Background(), gomock.Any()).
		Return(user, nil)

	service := NewUserService(repository, nil, log)

	u, err := service.GetUser(context.Background(), user.ID())
	assert.Equal(t, user, u, "should return user")
	assert.NoError(t, err, "should not error")
}
