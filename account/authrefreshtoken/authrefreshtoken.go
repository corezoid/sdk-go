package authrefreshtoken

import (
	"encoding/json"
	"fmt"
	"github.com/corezoid/sdk-go/account/oauth"
	"net/http"
	"net/url"
	"strings"
)

type Result struct {
	Result             string `json:"result"`
	Description        string `json:"description"`
	UserId             int `json:"user_id"`
	NewAccessToken     oauth.AccessToken `json:"new_access_token"`
	NewAccessTokenExpireAt int64 `json:"new_access_token_expire"`
	RefreshToken       oauth.RefreshToken
	Err error
	Resp *http.Response
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

func (api *Api) Request(token oauth.RefreshToken) *Result {
	result := &Result{RefreshToken: token}

	form := url.Values{}
	form.Set("client_id", api.c.ClientId)
	form.Set("grant_type", "refresh_token")
	form.Set("client_secret", api.c.ClientSecret)
	form.Set("refresh_token", string(token))

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/oauth2/token", api.c.HttpHost),
		strings.NewReader(form.Encode()),
	)

	if err != nil {
		result.Err = err

		return result
	}

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
