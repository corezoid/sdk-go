package accesstokenscopes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/corezoid/sdk-go/account/oauth"
)

type Result struct {
	Result string  `json:"result"`
	Scopes []Scope `json:"scopes"`
	Err    error
	Resp   *http.Response
}

type Scope struct {
	Scope  string `json:"scope"`
	Result string `json:"result"`
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

func (api *Api) Request(token oauth.AccessToken, scopes []string) *Result {
	result := &Result{}

	form := url.Values{}
	form.Set("client_id", api.c.ClientId)
	form.Set("client_secret", api.c.ClientSecret)
	form.Set("access_token", string(token))
	form.Set("scope", strings.Join(scopes, ","))

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/oauth2/scope", api.c.HttpHost),
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

	err := json.NewDecoder(r.Resp.Body).Decode(r)
	if err != nil {
		r.Err = err

		return r
	}

	if r.Result != "ok" {
		r.Err = fmt.Errorf("result not ok: %s", r.Result)

		return r
	}

	return r
}

func (r *Result) Close() {
	if r.Resp != nil {
		r.Resp.Body.Close()
	}
}
