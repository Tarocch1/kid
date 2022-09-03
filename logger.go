package kid

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type LoggerLevel string

const (
	LoggerLevelInfo  LoggerLevel = "info"
	LoggerLevelError LoggerLevel = "error"
)

const maxExtraLength = 32 << 10 // 32 KB

type Logger struct {
	module string
}

func NewLogger(module string) *Logger {
	return &Logger{
		module: module,
	}
}

func (l *Logger) log(data string) {
	fmt.Println(data)
}

func (l *Logger) FormatMessage(
	c *Ctx,
	level LoggerLevel,
	message string,
	extra map[string]interface{},
	err error,
) string {
	var extraStr string
	if extra != nil {
		extraBytes, _err := json.Marshal(extra)
		if _err != nil {
			logger.Error(c, "logger format message error", nil, _err)
		} else if len(extraBytes) > maxExtraLength {
			extraStr = "too long to show"
		} else {
			extraStr = string(extraBytes)
		}
	}

	var requestId string
	if c != nil {
		if id := c.Get(CtxRequestId); id != nil {
			requestId = id.(string)
		}
	}

	units := []string{
		time.Now().Format(time.RFC3339),
		fmt.Sprintf("| %s", level),
		_if(requestId != "", fmt.Sprintf("| %s", requestId), ""),
		_if(l.module != "", fmt.Sprintf("| [%s]", l.module), ""),
		message,
		_if(extraStr != "", fmt.Sprintf("- Extra %s", extraStr), ""),
		_if(err != nil, fmt.Sprintf("- Error %s", err), ""),
	}
	units = sliceFilter(units, func(item string) bool {
		return item != ""
	})
	return strings.ReplaceAll((strings.Join(units, " ")), "\n", "")
}

func (l *Logger) Info(c *Ctx, message string, extra map[string]interface{}) {
	l.log(l.FormatMessage(c, LoggerLevelInfo, message, extra, nil))
}

func (l *Logger) Error(c *Ctx, message string, extra map[string]interface{}, err error) {
	l.log(l.FormatMessage(c, LoggerLevelError, message, extra, err))
}

var logger = NewLogger("Logger")
