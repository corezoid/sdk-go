package corezoid

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type auth interface {
	sign(req *http.Request) error
}

type apiKeyAuth struct {
	login  int
	secret string
}

type saTokenAuth struct {
	token        string
	encodedToken string
}

func (a *apiKeyAuth) sign(req *http.Request) error {
	payload, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	timestamp := time.Now().Unix()

	var signature string
	if strings.Contains(req.Header.Get("Content-Type"), "multipart/form-data") {
		signature = a.genMultipartSignature(payload, timestamp)
	} else {
		signature = a.genSignature(payload, timestamp)
	}

	req.URL.Path = fmt.Sprintf(
		"%s/%d/%d/%s",
		req.URL.Path,
		a.login,
		timestamp,
		signature,
	)

	req.Body = ioutil.NopCloser(bytes.NewReader(payload))

	return nil
}

func (a *saTokenAuth) sign(req *http.Request) error {
	if a.encodedToken == "" {
		a.encodedToken = base64.StdEncoding.EncodeToString([]byte(a.token))
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.encodedToken))

	return nil
}

func (a *apiKeyAuth) genMultipartSignature(payload []byte, timestamp int64) string {
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

	return a.genSignature([]byte(result), timestamp)
}

func (a *apiKeyAuth) genSignature(payload []byte, timestamp int64) string {
	sha := sha1.Sum([]byte(fmt.Sprintf("%d%s%s%s", timestamp, a.secret, payload, a.secret)))

	return hex.EncodeToString(sha[:])
}
