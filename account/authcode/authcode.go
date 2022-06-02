package authcode

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/corezoid/sdk-go/account/oauth"
)

type Result struct {
	UserId               int64              `json:"user_id"`
	Token                oauth.AccessToken  `json:"access_token"`
	ExpireAt             int64              `json:"access_token_expire"`
	RefreshToken         oauth.RefreshToken `json:"refresh_token"`
	RefreshTokenExpireAt int64              `json:"refresh_token_expire"`
	Result               string             `json:"result"`
	Description          string             `json:"description"`
	Code                 oauth.AuthCode
	Err                  error
	Resp                 *http.Response
}

type Api struct {
	c oauth.Client
	h *http.Client
}

func New(c oauth.Client, h *http.Client) *Api {
	if h == nil {
		h = http.DefaultClient
	}

	return &Api{c: c, h: h}
}

func (api *Api) Request(code oauth.AuthCode) *Result {
	result := &Result{Code: code}

	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", api.c.ClientId)
	form.Set("client_secret", api.c.ClientSecret)
	form.Set("redirect_uri", api.c.RedirectUri)
	form.Set("code", string(code))

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/oauth2/token", api.c.HttpHost),
		strings.NewReader(form.Encode()),
	)

	if err != nil {
		result.Err = err

		return result
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := api.h.Do(req)
	if err != nil {
		result.Err = err

		return result
	}

	result.Resp = resp

	return result
}

func (r *Result) Decode() *Result {
	if r.Err != nil {
		return r
	}

	if r.Resp.StatusCode != http.StatusOK {
		r.Err = fmt.Errorf("single account response not OK, got %d", r.Resp.StatusCode)

		return r
	}

	err := json.NewDecoder(r.Resp.Body).Decode(&r)
	if err != nil {
		r.Err = err

		return r
	}

	if r.Result != "ok" {
		r.Err = fmt.Errorf("result not ok: %s - %s", r.Result, r.Description)

		return r
	}

	return r
}

func (r *Result) Close() {
	if r.Resp != nil {
		r.Resp.Body.Close()
	}
}
