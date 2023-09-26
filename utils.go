package goapi

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"strconv"

	"github.com/NaturalSelectionLabs/jschema"
)

func fnName(fn interface{}) string {
	fv := reflect.ValueOf(fn)

	fi := runtime.FuncForPC(fv.Pointer())

	return toPathName(regexp.MustCompile(`^.+\.`).ReplaceAllString(fi.Name(), ""))
}

var tUnmarshaler = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()

// converts the val to the kind of value.
func toValue(t reflect.Type, val string) (reflect.Value, error) {
	if t.Kind() == reflect.String || t.Implements(tUnmarshaler) || reflect.New(t).Type().Implements(tUnmarshaler) {
		val = strconv.Quote(val)
	}

	v := reflect.New(t)

	err := json.Unmarshal([]byte(val), v.Interface())
	if err != nil {
		return reflect.Value{}, fmt.Errorf("can't parse `%s` to expected value, %w", val, err)
	}

	return v.Elem(), nil
}

func tagName(t reflect.StructTag, name string) string {
	tag := jschema.ParseJSONTag(t)

	if tag != nil && tag.Name != "" {
		return tag.Name
	}

	return name
}

func firstProp(s *jschema.Schema) (p *jschema.Schema) { //nolint: nonamedreturns
	for _, p = range s.Properties {
		break
	}

	return p
}
