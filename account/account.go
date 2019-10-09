package account

import (
	"github.com/corezoid/sdk-go/account/authcode"
	"github.com/corezoid/sdk-go/account/authrefreshtoken"
	"github.com/corezoid/sdk-go/account/userinfo"
	"net/http"
	"github.com/corezoid/sdk-go/account/oauth"
)

type Service struct {
	authCodeApi *authcode.Api
	authRefreshTokenApi *authrefreshtoken.Api
	userInfoApi *userinfo.Api
}

func New(c oauth.Client, h *http.Client) *Service {
	if h == nil {
		h = http.DefaultClient
	}

	return &Service{
		userInfoApi: userinfo.New(c, h),
		authCodeApi: authcode.New(c, h),
		authRefreshTokenApi: authrefreshtoken.New(c, h),
	}
}

func (s *Service) UserInfo(token oauth.AccessToken) *userinfo.Result {
	return s.userInfoApi.Request(token)
}

func (s *Service) AuthCode(code oauth.AuthCode) *authcode.Result {
	return s.authCodeApi.Request(code)
}

func (s *Service) AuthRefreshToken(token oauth.RefreshToken) *authrefreshtoken.Result {
	return s.authRefreshTokenApi.Request(token)
}
