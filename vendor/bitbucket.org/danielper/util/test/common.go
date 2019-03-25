package test

import (
	"encoding/json"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo"
)

var e = echo.New()

func setupReq(content interface{}, endpoint, method string) (*httptest.ResponseRecorder, echo.Context) {
	cnt, _ := json.Marshal(content)

	req := httptest.NewRequest(method, endpoint, strings.NewReader(string(cnt)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return rec, c
}
