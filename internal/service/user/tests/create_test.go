package tests

import (
	"context"
	"testing"

	txManagerMocks "github.com/WithSoull/UserServer/internal/client/db/mocks"
	"github.com/WithSoull/UserServer/internal/model"
	"github.com/WithSoull/UserServer/internal/repository"
	userRepositoryMocks "github.com/WithSoull/UserServer/internal/repository/mocks"
	userService "github.com/WithSoull/UserServer/internal/service/user"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreate(t *testing.T) {
	type userRepoMockFunc func(t *testing.T, mc *minimock.Controller) repository.UserRepository

	type args struct {
		ctx             context.Context
		userInfo        model.UserInfo
		password        string
		passwordConfirm string
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		name     = gofakeit.Name()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, true, true, true, true, 8)

		info = model.UserInfo{
			Name:  name,
			Email: email,
		}
		defaultUserRepositoryMockFunc = func(t *testing.T, mc *minimock.Controller) repository.UserRepository {
			mock := userRepositoryMocks.NewUserRepositoryMock(mc)
			return mock
		}
	)

	tests := []struct {
		name         string
		args         args
		want_id      int64
		want_code    codes.Code
		err          error
		userRepoMock userRepoMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:             ctx,
				userInfo:        info,
				password:        password,
				passwordConfirm: password,
			},
			want_id:   id,
			want_code: codes.OK,
			err:       nil,
			userRepoMock: func(t *testing.T, mc *minimock.Controller) repository.UserRepository {
				mock := userRepositoryMocks.NewUserRepositoryMock(mc)
				mock.CreateMock.Inspect(func(ctx context.Context, user *model.UserInfo, hashedPassword string) {
					require.NotEmpty(t, hashedPassword)
				}).Return(id, nil)
				return mock
			},
		},
		{
			name: "name is empty case",
			args: args{
				ctx: ctx,
				userInfo: model.UserInfo{
					Name:  "",
					Email: email,
				},
				password:        password,
				passwordConfirm: password,
			},
			want_id:      id,
			want_code:    codes.InvalidArgument,
			err:          nil,
			userRepoMock: defaultUserRepositoryMockFunc,
		},
		{
			name: "email is empty case",
			args: args{
				ctx: ctx,
				userInfo: model.UserInfo{
					Name:  name,
					Email: "",
				},
				password:        password,
				passwordConfirm: password,
			},
			want_id:      id,
			want_code:    codes.InvalidArgument,
			err:          nil,
			userRepoMock: defaultUserRepositoryMockFunc,
		},
		{
			name: "invalid mail case",
			args: args{
				ctx: ctx,
				userInfo: model.UserInfo{
					Name:  name,
					Email: "ivalid_mail@ru",
				},
				password:        password,
				passwordConfirm: password,
			},
			want_id:      id,
			want_code:    codes.InvalidArgument,
			err:          nil,
			userRepoMock: defaultUserRepositoryMockFunc,
		},
		{
			name: "password is empty case",
			args: args{
				ctx:             ctx,
				userInfo:        info,
				password:        "",
				passwordConfirm: "",
			},
			want_id:      id,
			want_code:    codes.InvalidArgument,
			err:          nil,
			userRepoMock: defaultUserRepositoryMockFunc,
		},
		{
			name: "password is not queal case",
			args: args{
				ctx:             ctx,
				userInfo:        info,
				password:        "12345",
				passwordConfirm: "54321",
			},
			want_id:      id,
			want_code:    codes.InvalidArgument,
			err:          nil,
			userRepoMock: defaultUserRepositoryMockFunc,
		},
		{
			name: "password is too short case",
			args: args{
				ctx:             ctx,
				userInfo:        info,
				password:        "1234",
				passwordConfirm: "1234",
			},
			want_id:      id,
			want_code:    codes.InvalidArgument,
			err:          nil,
			userRepoMock: defaultUserRepositoryMockFunc,
		},
	}

	for _, tt := range tests {
		tt := tt // To avoid bugs in parralel tests
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userRepoMock := tt.userRepoMock(t, mc)
			txManagerMock := txManagerMocks.NewTxManagerMock(mc)

			service := userService.NewService(userRepoMock, txManagerMock)

			res_id, err := service.Create(ctx, tt.args.userInfo, tt.args.password, tt.args.passwordConfirm)
			if tt.want_code == codes.OK {
				require.NoError(t, err)
				require.Equal(t, res_id, id)
			} else {
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.want_code, st.Code())
			}
		})
	}
}
