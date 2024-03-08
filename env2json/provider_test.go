package env2json

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnv_Read(t *testing.T) {
	prov := Provider("", "", nil)

	_, err := prov.Read()

	require.Error(t, err)
}

func TestEnv_ReadBytes(t *testing.T) {
	type fields struct {
		prefix   string
		delim    string
		callback func(key string, value string) (string, interface{})
		out      string
	}

	tests := []struct {
		name    string
		fields  fields
		envVars map[string]string
		want    []byte
		wantErr bool
	}{
		{
			name: "some data",
			fields: fields{
				prefix:   "UNITTESTPREFIX",
				delim:    "_",
				callback: nil,
				out:      "{}",
			},
			envVars: map[string]string{
				"UNITTESTPREFIX_SOME": "data",
			},
			want:    []byte(`{"UNITTESTPREFIX":{"SOME":"data"}}`),
			wantErr: false,
		},
		{
			name: "drop key",
			fields: fields{
				prefix: "UNITTESTPREFIX",
				delim:  "_",
				callback: func(key string, value string) (string, interface{}) {
					if strings.Contains(key, "DROP") {
						return "", value
					}

					return strings.ToLower(strings.TrimPrefix(key, "UNITTESTPREFIX_")), value
				},
				out: "{}",
			},
			envVars: map[string]string{
				"UNITTESTPREFIX_SOME": "data",
				"UNITTESTPREFIX_DROP": "data",
			},
			want:    []byte(`{"some":"data"}`),
			wantErr: false,
		},
		{
			name: "empty prefix",
			fields: fields{
				prefix: "",
				delim:  "_",
				callback: func(key string, value string) (string, interface{}) {
					if !strings.Contains(key, "UNITTESTPREFIX") {
						return "", value
					}

					return strings.ToLower(strings.TrimPrefix(key, "UNITTESTPREFIX_")), value
				},
				out: "{}",
			},
			envVars: map[string]string{
				"UNITTESTPREFIX_SOME": "data",
			},
			want:    []byte(`{"some":"data"}`),
			wantErr: false,
		},
		{
			name: "some data with callback",
			fields: fields{
				prefix: "UNITTESTPREFIX",
				delim:  "_",
				callback: func(key string, value string) (string, interface{}) {
					return strings.ToLower(strings.TrimPrefix(key, "UNITTESTPREFIX_")), value
				},
				out: "{}",
			},
			envVars: map[string]string{
				"UNITTESTPREFIX_SOME": "data",
			},
			want:    []byte(`{"some":"data"}`),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Env{
				prefix:   tt.fields.prefix,
				delim:    tt.fields.delim,
				callback: tt.fields.callback,
				out:      tt.fields.out,
			}

			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			got, err := e.ReadBytes()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadBytes() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadBytes() got = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestProvider(t *testing.T) {
	type args struct {
		prefix   string
		delim    string
		callback func(s string) string
	}

	tests := []struct {
		name string
		args args
		want *Env
	}{
		{
			name: "default",
			args: args{
				prefix:   "test",
				delim:    "_",
				callback: nil,
			},
			want: &Env{
				prefix:   "test",
				delim:    "_",
				callback: nil,
				out:      "{}",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Provider(tt.args.prefix, tt.args.delim, tt.args.callback); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Provider() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProvider_Callback(t *testing.T) {
	type args struct {
		prefix   string
		delim    string
		callback func(s string) string
	}

	tests := []struct {
		name string
		args args
		want *Env
	}{
		{
			name: "withCallback",
			args: args{
				prefix: "test",
				delim:  "_",
				callback: func(s string) string {
					return s
				},
			},
			want: &Env{
				prefix: "test",
				delim:  "_",
				callback: func(key string, value string) (string, interface{}) {
					return func(s string) string {
						return s
					}(key), value
				},
				out: "{}",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := Provider(tt.args.prefix, tt.args.delim, tt.args.callback)

			gotKey, gotVal := provider.callback("some", "data")
			wantKey, wantVal := tt.want.callback("some", "data")

			if !reflect.DeepEqual(gotKey, wantKey) ||
				!reflect.DeepEqual(gotVal, wantVal) {
				t.Errorf("Provider() = %v=%v, want %v=%v", gotKey, gotVal, wantKey, wantVal)
			}
		})
	}
}

func TestProviderWithValue(t *testing.T) {
	type args struct {
		prefix   string
		delim    string
		callback func(key string, value string) (string, interface{})
	}

	tests := []struct {
		name string
		args args
		want *Env
	}{
		{
			name: "default",
			args: args{
				prefix:   "test",
				delim:    "_",
				callback: nil,
			},
			want: &Env{
				prefix:   "test",
				delim:    "_",
				callback: nil,
				out:      "{}",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProviderWithValue(tt.args.prefix, tt.args.delim, tt.args.callback); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProviderWithValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProviderWithValue_Callback(t *testing.T) {
	type args struct {
		prefix   string
		delim    string
		callback func(key string, value string) (string, interface{})
	}

	tests := []struct {
		name string
		args args
		want *Env
	}{
		{
			name: "withCallback",
			args: args{
				prefix: "test",
				delim:  "_",
				callback: func(key string, value string) (string, interface{}) {
					return key, value
				},
			},
			want: &Env{
				prefix: "test",
				delim:  "_",
				callback: func(key string, value string) (string, interface{}) {
					return key, value
				},
				out: "{}",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := ProviderWithValue(tt.args.prefix, tt.args.delim, tt.args.callback)
			gotKey, gotVal := provider.callback("some", "data")
			wantKey, wantVal := tt.want.callback("some", "data")

			if !reflect.DeepEqual(gotKey, wantKey) ||
				!reflect.DeepEqual(gotVal, wantVal) {
				t.Errorf("ProviderWithValue() = %v=%v, want %v=%v", gotKey, gotVal, wantKey, wantVal)
			}
		})
	}
}
