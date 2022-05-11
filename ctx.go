package kid

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Ctx struct {
	writer  http.ResponseWriter
	request *http.Request

	params  map[string]string
	rawBody []byte

	status   int
	store    map[string]interface{}
	handlers []HandlerFunc
	index    int
}

func newCtx(w http.ResponseWriter, r *http.Request) *Ctx {
	rawBody, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewReader(rawBody))

	context := &Ctx{
		writer:  w,
		request: r,

		params:  make(map[string]string),
		rawBody: rawBody,

		status:   http.StatusOK,
		store:    make(map[string]interface{}),
		handlers: make([]HandlerFunc, 0),
		index:    -1,
	}
	return context
}

func (c *Ctx) Next() error {
	c.index++
	if c.index < len(c.handlers) {
		handler := c.handlers[c.index]
		return handler(c)
	}
	return nil
}

func (c *Ctx) Set(key string, value interface{}) {
	c.store[key] = value
}

func (c *Ctx) Get(key string) interface{} {
	return c.store[key]
}

func (c *Ctx) Method() string {
	return c.request.Method
}

func (c *Ctx) Url() *url.URL {
	return c.request.URL
}

func (c *Ctx) Params() map[string]string {
	return c.params
}

func (c *Ctx) GetParam(key string, defaultValue ...string) string {
	return getValue(c.params[key], defaultValue...)
}

func (c *Ctx) Query() url.Values {
	return c.request.URL.Query()
}

func (c *Ctx) GetQuery(key string, defaultValue ...string) string {
	return getValue(c.request.URL.Query().Get(key), defaultValue...)
}

func (c *Ctx) Header() http.Header {
	return c.request.Header
}

func (c *Ctx) GetHeader(key string, defaultValue ...string) string {
	return getValue(c.request.Header.Get(key), defaultValue...)
}

func (c *Ctx) GetHeaderValues(key string, defaultValue ...[]string) []string {
	return getValue(c.request.Header.Values(key), defaultValue...)
}

func (c *Ctx) Cookies() []*http.Cookie {
	return c.request.Cookies()
}

func (c *Ctx) GetCookie(name string) *http.Cookie {
	cookie, _ := c.request.Cookie(name)
	return cookie
}

func (c *Ctx) FormValue(key string, defaultValue ...string) string {
	return getValue(c.request.FormValue(key), defaultValue...)
}

func (c *Ctx) FormFile(key string) (*multipart.FileHeader, error) {
	_, fh, err := c.request.FormFile(key)
	return fh, err
}

func (c *Ctx) Body() []byte {
	return c.rawBody
}

func (c *Ctx) BodyParser(out interface{}) error {
	ctype := c.GetHeader("Content-Type")

	switch {
	case strings.HasPrefix(ctype, "application/json"):
		return json.Unmarshal(c.Body(), out)
	case strings.HasPrefix(ctype, "application/x-www-form-urlencoded") ||
		c.Method() == "HEAD" ||
		c.Method() == "GET" ||
		c.Method() == "DELETE":
		err := c.request.ParseForm()
		if err != nil {
			return err
		}
		return unmarshalForm(c.request.Form, out)
	case strings.HasPrefix(ctype, "multipart/form-data"):
		err := c.request.ParseMultipartForm(32 << 20)
		if err != nil {
			return err
		}
		return unmarshalForm(c.request.MultipartForm.Value, out)
	case strings.HasPrefix(ctype, "text/xml") ||
		strings.HasPrefix(ctype, "application/xml"):
		return xml.Unmarshal(c.Body(), out)
	}

	return NewError(http.StatusUnprocessableEntity, "422 Unprocessable Entity")
}

func (c *Ctx) SetHeader(key string, value string) *Ctx {
	c.writer.Header().Set(key, value)
	return c
}

func (c *Ctx) AddHeader(key string, value string) *Ctx {
	c.writer.Header().Add(key, value)
	return c
}

func (c *Ctx) SetCookie(cookie *http.Cookie) *Ctx {
	http.SetCookie(c.writer, cookie)
	return c
}

func (c *Ctx) ClearCookie(name string) *Ctx {
	http.SetCookie(c.writer, &http.Cookie{
		Name:     name,
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	})
	return c
}

func (c *Ctx) Status(status int) *Ctx {
	c.status = status
	return c
}

func (c *Ctx) SendStatus(status int) error {
	return c.Status(status).String(strconv.Itoa(status))
}

func (c *Ctx) Data(data []byte) error {
	c.SetHeader("Content-Type", "application/octet-stream")
	c.writer.WriteHeader(c.status)
	_, err := c.writer.Write(data)
	return err
}

func (c *Ctx) String(format string, values ...interface{}) error {
	c.SetHeader("Content-Type", "text/plain")
	c.writer.WriteHeader(c.status)
	_, err := c.writer.Write([]byte(fmt.Sprintf(format, values...)))
	return err
}

func (c *Ctx) JSON(data interface{}) error {
	c.SetHeader("Content-Type", "application/json")
	c.writer.WriteHeader(c.status)
	encoder := json.NewEncoder(c.writer)
	return encoder.Encode(data)
}

func (c *Ctx) HTML(code int, html string) error {
	c.SetHeader("Content-Type", "text/html")
	c.writer.WriteHeader(c.status)
	_, err := c.writer.Write([]byte(html))
	return err
}
