package corezoid

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Result struct {
	RequestProc string  `json:"request_proc"`
	Ops         []MapOp `json:"ops"`
	Response    *http.Response
	Err         error
}

func (r *Result) Decode() *Result {
	if r.Err != nil {
		return r
	}
	if r.Response == nil {
		r.Err = errors.New("response is nil")

		return r
	}
	if r.Response.StatusCode != http.StatusOK {
		r.Err = fmt.Errorf("response status not ok, but %d", r.Response.StatusCode)

		return r
	}

	if err := json.NewDecoder(r.Response.Body).Decode(&r); err != nil {
		r.Err = err
		r.Response.Body.Close()
	}

	if r.RequestProc != "ok" {
		r.Err = fmt.Errorf("request_proc not ok")
	}

	return r
}

func (r *Result) Close() error {
	if r.Response != nil {
		return r.Response.Body.Close()
	}

	return nil
}
