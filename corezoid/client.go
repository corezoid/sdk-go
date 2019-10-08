package corezoid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/textproto"
	"strconv"
	"strings"
)

type Client struct {
	Endpoint   string
	HttpClient *http.Client
}

// NewApiKey creates a client that uses api key to access API
func New(endpoint string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		Endpoint: endpoint,
		HttpClient: httpClient,
	}
}

func NewCloud(httpClient *http.Client) *Client {
	return New("https://api.corezoid.com", httpClient)
}

func (c *Client) Call(ops Ops, auth Auth) (result *Result) {
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

	err = auth.Sign(req)
	if err != nil {
		result.Err = err

		return result
	}

	b, _ := httputil.DumpRequest(req, true)
	log.Print("REQUEST: ", string(b))

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		result.Err = err

		return result
	}
	result.Response = resp

	b, _ = httputil.DumpResponse(resp, true)
	log.Print("RESPONSE: ", string(b))

	return result
}

func (c *Client) Upload(op Op, auth Auth) (result *Result) {
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
		var err error
		switch t := val.(type) {
		case int64:
			err = writer.WriteField(key, strconv.FormatInt(val.(int64), 10))
		case string:
			err = writer.WriteField(key, val.(string))
		default:
			err = fmt.Errorf("key %s has unsupported type: %s", key, t)
		}

		if err != nil {
			result.Err = err

			return result
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

	err = auth.Sign(req)
	if err != nil {
		result.Err = err

		return result
	}

	b, _ := httputil.DumpRequest(req, true)
	log.Print("REQUEST: ", string(b))

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		result.Err = err

		return result
	}

	b, _ = httputil.DumpResponse(resp, true)
	log.Print("RESPONSE: ", string(b))

	result.Response = resp

	return result
}
