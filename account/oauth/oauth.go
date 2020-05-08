package oauth

type Client struct {
	ClientId     string
	ClientSecret string
	HttpHost     string
	RedirectUri  string
}

type AccessToken string
type RefreshToken string
type AuthCode string
