package stringtomap

import (
	"reflect"

	"github.com/go-viper/mapstructure/v2"
	"gopkg.in/yaml.v3"
)

// StringToMap converts a string to a map
func StringToMap() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t.Kind() != reflect.Map {
			return data, nil
		}

		if mapString, ok := (data).(string); ok {
			out := reflect.New(t).Interface()

			if err := yaml.Unmarshal([]byte(mapString), &out); err != nil {
				return data, err
			}
			return out, nil
		}

		return data, nil
	}
}
