package env2json

import (
	"errors"
	"os"
	"strings"

	"github.com/tidwall/sjson"
)

type Env struct {
	prefix   string
	delim    string
	callback func(key string, value string) (string, interface{})
	out      string
}

var errUnsupportedMethod = errors.New("env2json provider does not support this method")

func Provider(prefix, delim string, callback func(s string) string) *Env {
	env := &Env{
		prefix: prefix,
		delim:  delim,
		out:    "{}",
	}
	if callback != nil {
		env.callback = func(key string, value string) (string, interface{}) {
			return callback(key), value
		}
	}

	return env
}

func ProviderWithValue(prefix, delim string, cb func(key string, value string) (string, interface{})) *Env {
	return &Env{
		prefix:   prefix,
		delim:    delim,
		out:      "{}",
		callback: cb,
	}
}

// ReadBytes reads the contents of a file on disk and returns the bytes.
func (e *Env) ReadBytes() ([]byte, error) {
	// Collect the environment variable keys.
	var keys []string

	for _, key := range os.Environ() {
		if e.prefix != "" {
			if strings.HasPrefix(key, e.prefix) {
				keys = append(keys, key)
			}
		} else {
			keys = append(keys, key)
		}
	}

	for _, k := range keys {
		//nolint:gomnd
		parts := strings.SplitN(k, "=", 2)

		var (
			key   string
			value interface{}
		)

		// If there's a transformation callback,
		// run it through every key/value.
		if e.callback != nil {
			key, value = e.callback(parts[0], parts[1])
			// If the callback blanked the key, it should be omitted
			if key == "" {
				continue
			}
		} else {
			key = parts[0]
			value = parts[1]
		}

		if err := e.set(key, value); err != nil {
			return []byte{}, err
		}
	}

	return []byte(e.out), nil
}

func (e *Env) set(key string, value interface{}) error {
	out, err := sjson.Set(e.out, strings.ReplaceAll(key, e.delim, "."), value)
	if err != nil {
		return err
	}

	e.out = out

	return nil
}

// Read is not supported by the file provider.
func (e *Env) Read() (map[string]interface{}, error) {
	return nil, errUnsupportedMethod
}
