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
	"os"
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

// Next starts the next middleware.
func (c *Ctx) Next() error {
	c.index++
	if c.index < len(c.handlers) {
		handler := c.handlers[c.index]
		return handler(c)
	}
	return nil
}

// Set sets some value to ctx.
func (c *Ctx) Set(key string, value interface{}) {
	c.store[key] = value
}

// Get gets some value from ctx.
func (c *Ctx) Get(key string) interface{} {
	return c.store[key]
}

// Method returns request's method.
func (c *Ctx) Method() string {
	return c.request.Method
}

// Url returns request's URL.
func (c *Ctx) Url() *url.URL {
	return c.request.URL
}

// Params gets all router path params.
func (c *Ctx) Params() map[string]string {
	return c.params
}

// GetParam gets a router path param value by key.
func (c *Ctx) GetParam(key string, defaultValue ...string) string {
	return getValue(c.params[key], defaultValue...)
}

// Query gets request's Query.
func (c *Ctx) Query() url.Values {
	return c.request.URL.Query()
}

// GetQuery gets a query value by key.
func (c *Ctx) GetQuery(key string, defaultValue ...string) string {
	return getValue(c.request.URL.Query().Get(key), defaultValue...)
}

// Header gets request's Header.
func (c *Ctx) Header() http.Header {
	return c.request.Header
}

// GetHeader gets a header's first value  by key.
func (c *Ctx) GetHeader(key string, defaultValue ...string) string {
	return getValue(c.request.Header.Get(key), defaultValue...)
}

// GetHeaderValues gets a header by key.
func (c *Ctx) GetHeaderValues(key string, defaultValue ...[]string) []string {
	return getValue(c.request.Header.Values(key), defaultValue...)
}

// Cookies gets request's cookies.
func (c *Ctx) Cookies() []*http.Cookie {
	return c.request.Cookies()
}

// GetCookie gets a cookie by name.
func (c *Ctx) GetCookie(name string) *http.Cookie {
	cookie, _ := c.request.Cookie(name)
	return cookie
}

// FormValue gets a form value by key.
func (c *Ctx) FormValue(key string, defaultValue ...string) string {
	return getValue(c.request.FormValue(key), defaultValue...)
}

// FormFile gets a form file by key.
func (c *Ctx) FormFile(key string) (*multipart.FileHeader, error) {
	_, fh, err := c.request.FormFile(key)
	return fh, err
}

// Body gets request's raw body.
func (c *Ctx) Body() []byte {
	return c.rawBody
}

// BodyParser parsers body to any struct.
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

// SetHeader sets a header.
func (c *Ctx) SetHeader(key string, value string) *Ctx {
	c.writer.Header().Set(key, value)
	return c
}

// AddHeader adds a header value.
func (c *Ctx) AddHeader(key string, value string) *Ctx {
	c.writer.Header().Add(key, value)
	return c
}

// SetCookie sets a cookie.
func (c *Ctx) SetCookie(cookie *http.Cookie) *Ctx {
	http.SetCookie(c.writer, cookie)
	return c
}

// ClearCookie clears a cookie.
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

// Status sets response's status.
func (c *Ctx) Status(status int) *Ctx {
	c.status = status
	return c
}

// SendStatus sets response's status and send.
func (c *Ctx) SendStatus(status int) error {
	return c.Status(status).String(strconv.Itoa(status))
}

// Redirect redirects request to target with status.
func (c *Ctx) Redirect(target string, status ...int) error {
	_status := http.StatusTemporaryRedirect
	if len(status) > 0 {
		_status = status[0]
	}
	c.SetHeader("Location", target)
	return c.SendStatus(_status)
}

// Stream sends binary stream.
func (c *Ctx) Stream(data []byte) error {
	c.SetHeader("Content-Type", "application/octet-stream; charset=utf-8")
	c.writer.WriteHeader(c.status)
	_, err := c.writer.Write(data)
	return err
}

// String sends string.
func (c *Ctx) String(format string, values ...interface{}) error {
	c.SetHeader("Content-Type", "text/plain; charset=utf-8")
	c.writer.WriteHeader(c.status)
	_, err := c.writer.Write([]byte(fmt.Sprintf(format, values...)))
	return err
}

// Json sends json.
func (c *Ctx) Json(data interface{}) error {
	c.SetHeader("Content-Type", "application/json; charset=utf-8")
	c.writer.WriteHeader(c.status)
	encoder := json.NewEncoder(c.writer)
	return encoder.Encode(data)
}

// Html sends html.
func (c *Ctx) Html(html string) error {
	c.SetHeader("Content-Type", "text/html; charset=utf-8")
	c.writer.WriteHeader(c.status)
	_, err := c.writer.Write([]byte(html))
	return err
}

// SendFile reads file from fs at path and sends it.
func (c *Ctx) SendFile(path string, download bool, fs ...http.FileSystem) error {
	var _fs http.FileSystem
	if len(fs) > 0 {
		_fs = fs[0]
	} else {
		dir, err := os.Getwd()
		if err != nil {
			return nil
		}
		_fs = http.Dir(dir)
	}

	file, err := _fs.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	if stat.IsDir() {
		return NewError(http.StatusBadRequest, "400 Bad Request: Can not serve dir")
	}

	if download {
		c.SetHeader("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", stat.Name()))
	}

	http.ServeContent(c.writer, c.request, stat.Name(), stat.ModTime(), file)
	return nil
}
