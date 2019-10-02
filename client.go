package corezoid

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"regexp"
	"strings"
	"time"
)

type Client struct {
	Endpoint   string
	HttpClient *http.Client
	secret     string
	login      int
}

func New(login int, secret string) *Client {
	return &Client{
		Endpoint:   "https://api.corezoid.com",
		HttpClient: http.DefaultClient,
		login:      login,
		secret:     secret,
	}
}

func (c *Client) Call(ops Ops) (result *Result) {
	result = &Result{}

	payload, err := json.Marshal(ops.Raw())
	if err != nil {
		result.Err = err

		return result
	}

	timestamp := time.Now().Unix()

	uri := fmt.Sprintf(
		"%s/api/2/json/%d/%d/%s",
		c.Endpoint,
		c.login,
		timestamp,
		c.genSignature(payload, timestamp),
	)

	req, err := http.NewRequest("POST", uri, bytes.NewReader(payload))
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

	payload, err := ioutil.ReadAll(body)
	if err != nil {
		result.Err = err

		return result
	}

	timestamp := time.Now().Unix()

	uri := fmt.Sprintf(
		"%s/api/2/upload/%d/%d/%s",
		c.Endpoint,
		c.login,
		timestamp,
		c.genMultipartSignature(payload, timestamp),
	)

	req, err := http.NewRequest("POST", uri, bytes.NewReader(payload))
	if err != nil {
		result.Err = err

		return result
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		result.Err = err

		return result
	}

	result.Response = resp

	return result
}

func (c *Client) genMultipartSignature(payload []byte, timestamp int64) string {
	payload = regexp.MustCompile(`-+[\d\w]+-{0,}\r\n`).ReplaceAll(payload, []byte(""))
	payload = regexp.MustCompile(`\r\n\r\n`).ReplaceAll(payload, []byte("\r\n"))
	payload = []byte(strings.TrimSpace(string(payload)))

	chunks := strings.Split(string(payload), "\r\n")

	reg := regexp.MustCompile(`^Content-`)

	result := ""
	for _, chunk := range chunks {
		result = result + chunk

		if reg.Match([]byte(chunk)) {
			result = result + "\r\n"
		}
	}

	return c.genSignature([]byte(result), timestamp)
}

func (c *Client) genSignature(payload []byte, timestamp int64) string {
	sha := sha1.Sum([]byte(fmt.Sprintf("%d%s%s%s", timestamp, c.secret, payload, c.secret)))

	return hex.EncodeToString(sha[:])
}
