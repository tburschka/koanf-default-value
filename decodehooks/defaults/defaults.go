package defaults

import (
	"reflect"
	"strings"

	"github.com/go-viper/mapstructure/v2"
)

func Defaults(keyTag string, defaultTag string) mapstructure.DecodeHookFunc {
	return func(f reflect.Value, t reflect.Value) (interface{}, error) {
		tType := t.Type()

		if tType.Kind() == reflect.Struct {
			for i := 0; i < tType.NumField(); i++ {
				setDefault(f, tType.Field(i), keyTag, defaultTag)
			}
		}

		return f.Interface(), nil
	}
}

func setDefault(f reflect.Value, t reflect.StructField, keyTag string, defaultTag string) {
	if f.Kind() == reflect.Map {
		key, _, _ := strings.Cut(t.Tag.Get(keyTag), ",")
		val, ok := t.Tag.Lookup(defaultTag)
		if !ok {
			return
		}

		// check for key in the map
		for _, e := range f.MapKeys() {
			// key found, already set, nothing to do
			if key == e.String() {
				return
			}
		}

		fVal := reflect.ValueOf(val)

		// special handling for empty/missing structs
		tType := t.Type
		if tType.Kind() == reflect.Struct {
			// create an empty map
			fVal = reflect.ValueOf(map[string]any{})
			for i := 0; i < tType.NumField(); i++ {
				setDefault(fVal, tType.Field(i), keyTag, defaultTag)
			}
		}

		// add missing key with default to the map
		f.SetMapIndex(reflect.ValueOf(key), fVal)
	}
}
