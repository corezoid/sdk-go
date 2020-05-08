package authurl

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/corezoid/sdk-go/account/oauth"
)

type Api struct {
	c oauth.Client
}

func New(c oauth.Client) *Api {
	return &Api{c: c}
}

func (api *Api) AuthorizeUrl(scopes []string) string {
	return fmt.Sprintf(
		"%s/oauth2/authorize/?response_type=code&scope=%s&redirect_uri=%s&client_id=%s",
		api.c.HttpHost,
		strings.Join(scopes, ","),
		url.QueryEscape(api.c.RedirectUri),
		url.QueryEscape(api.c.ClientId),
	)
}
