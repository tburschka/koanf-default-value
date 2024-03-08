package main

import (
	"github.com/creasty/defaults"

	"config-debug/env2json"
	"github.com/go-viper/mapstructure/v2"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/v2"
	"github.com/subosito/gotenv"
	"strings"
)

const (
	configEnvPrefix = "APP"
	configEnvDelim  = "__"
)

type Config struct {
	String string   `yaml:"string,omitempty"`
	Bool   bool     `yaml:"bool,omitempty"`
	Struct Nested   `yaml:"struct,omitempty"`
	Slice  []Nested `yaml:"slice,omitempty" default:"[]"`
}

type Nested struct {
	TrueBool  bool `yaml:"true_bool,omitempty" default:"true"`
	FalseBool bool `yaml:"false_bool,omitempty" default:"false"`
}

func LoadConfig() *Config {
	_ = gotenv.Load()

	k := koanf.New(".")

	_ = k.Load(env2json.Provider(configEnvPrefix, configEnvDelim, func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, configEnvPrefix+configEnvDelim))
	}), json.Parser())

	example := &Config{}
	_ = k.UnmarshalWithConf("", example, koanf.UnmarshalConf{
		Tag: "yaml",
		DecoderConfig: &mapstructure.DecoderConfig{
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				mapstructure.StringToTimeDurationHookFunc(),
			),
			Metadata:         nil,
			Result:           example,
			WeaklyTypedInput: true,
		},
	})

	// set defaults
	_ = defaults.Set(example)

	return example
}
