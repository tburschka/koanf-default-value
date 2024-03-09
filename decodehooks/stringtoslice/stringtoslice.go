package stringtoslice

import (
	"reflect"

	"github.com/go-viper/mapstructure/v2"
	"gopkg.in/yaml.v3"
)

// StringToSlice converts a string to a slice
func StringToSlice() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t.Kind() != reflect.Slice {
			return data, nil
		}

		if sliceString, ok := (data).(string); ok {
			out := reflect.New(t).Interface()

			if err := yaml.Unmarshal([]byte(sliceString), &out); err != nil {
				return data, err
			}
			return out, nil
		}

		return data, nil
	}
}
