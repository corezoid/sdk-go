package userinfo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/corezoid/sdk-go/account/oauth"
)

type Result struct {
	UserId       int64  `json:"user_id"`
	Nick         string `json:"nick"`
	UserPhoto    string `json:"user_photo"`
	Login        string `json:"login"`
	Lang         string `json:"lang"`
	Status       string `json:"status"`
	CreationTime int    `json:"create_time"`
	Result       string `json:"result"`
	Description  string `json:"description"`
	Code         string
	Err          error
	Resp         *http.Response
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

func (api *Api) Request(token oauth.AccessToken) *Result {
	result := &Result{}

	form := url.Values{}
	form.Set("access_token", string(token))

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/oauth2/userinfo", api.c.HttpHost),
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

	err := json.NewDecoder(r.Resp.Body).Decode(r)
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
