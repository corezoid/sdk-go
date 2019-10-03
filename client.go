package corezoid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

type Client struct {
	Endpoint   string
	HttpClient *http.Client
	auth       auth
}

// NewApiKey creates a client that uses api key to access API
func NewApiKey(login int, secret string) *Client {
	return &Client{
		Endpoint:   "https://api.corezoid.com",
		HttpClient: http.DefaultClient,
		auth:       &apiKeyAuth{login: login, secret: secret},
	}
}

// NewSAToken creates a client that uses Single account (SA) access token to access API
func NewSAToken(token string) *Client {
	return &Client{
		Endpoint:   "https://api.corezoid.com",
		HttpClient: http.DefaultClient,
		auth:       &saTokenAuth{token: token},
	}
}

func (c *Client) Call(ops Ops) (result *Result) {
	result = &Result{}

	payload, err := json.Marshal(ops.Raw())
	if err != nil {
		result.Err = err

		return result
	}

	uri := fmt.Sprintf("%s/api/2/json", c.Endpoint)
	req, err := http.NewRequest("POST", uri, bytes.NewReader(payload))
	if err != nil {
		result.Err = err

		return result
	}
	req.Header.Set("Content-Type", "application/json")

	err = c.auth.sign(req)
	if err != nil {
		result.Err = err

		return result
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		result.Err = err

		return result
	}
	result.Response = resp

	return result
}

func (c *Client) Upload(op Op) (result *Result) {
	result = &Result{}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="scheme"; filename="scheme.json"`)
	h.Set("Content-Type", "application/json")
	part, err := writer.CreatePart(h)
	if err != nil {
		result.Err = err

		return result
	}

	rawOp := op.Raw()

	_, err = io.Copy(part, strings.NewReader(rawOp["scheme"].(string)))
	delete(rawOp, "scheme")

	for key, val := range rawOp {
		err := writer.WriteField(key, val.(string))
		if err != nil {
			panic(err)
		}
	}
	err = writer.Close()
	if err != nil {
		result.Err = err

		return result
	}

	uri := fmt.Sprintf("%s/api/2/upload", c.Endpoint)
	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		result.Err = err

		return result
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	err = c.auth.sign(req)
	if err != nil {
		result.Err = err

		return result
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		result.Err = err

		return result
	}

	result.Response = resp

	return result
}
