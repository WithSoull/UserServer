package tests

import (
	"context"
	"testing"

	"github.com/WithSoull/AuthService/internal/model"
	"github.com/WithSoull/AuthService/internal/service"
	desc "github.com/WithSoull/AuthService/pkg/user/v1"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
)

func TestCreate(t *testing.T) {
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req desc.CreateRequest
	}

	var (
		ctx = context.Background()
		mc = minimock.NewController()

		userInfo = model.UserInfo{
			Name: gofakeit.Name(),
			Email: gofakeit.Email(),
		}
	)

	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			}
		}
	}
}
