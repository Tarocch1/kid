package kid

import (
	"fmt"
	"net/url"
	"reflect"
)

func isZero[T interface{}](value T, defaultValue ...T) bool {
	v := reflect.ValueOf(value)
	return v.IsZero()
}

func getValue[T interface{}](value T, defaultValue ...T) T {
	if !isZero(value) {
		return value
	} else if len(defaultValue) > 0 {
		return defaultValue[0]
	} else {
		return value
	}
}

func unmarshalForm(data url.Values, out interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			if err, ok = r.(error); !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	outValue := reflect.ValueOf(out).Elem()
	outType := reflect.TypeOf(out).Elem()
	for i := 0; i < outValue.NumField(); i++ {
		field := outValue.Field(i)
		tag := outType.Field(i).Tag.Get("form")
		item, ok := data[tag]
		if ok {
			if field.Type().Kind() == reflect.String {
				field.SetString(item[0])
			} else {
				field.Set(reflect.ValueOf(item))
			}
		}
	}
	return nil
}
