package main

import (
	"config-debug/decodehooks/defaults"
	"config-debug/decodehooks/stringtomap"
	"config-debug/decodehooks/stringtoslice"
	"config-debug/providers/env2json"
	"dario.cat/mergo"
	"github.com/go-viper/mapstructure/v2"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/subosito/gotenv"
	"strings"
)

const (
	configEnvPrefix = "APP"
	configEnvDelim  = "__"
)

type Config struct {
	String  string   `yaml:"string,omitempty" default:"string"`
	Bool    bool     `yaml:"bool,omitempty"`
	Struct  Nested   `yaml:"struct,omitempty" default:"{}"`
	Slice   []Nested `yaml:"slice,omitempty" default:"[]"`
	Special Special  `yaml:"special,omitempty" default:"{}"`
}

type Nested struct {
	TrueBool  bool      `yaml:"true_bool,omitempty" default:"true"`
	FalseBool bool      `yaml:"false_bool,omitempty" default:"false"`
	SubBool   SubNested `yaml:"sub_bool,omitempty" default:"{}"`
}

type SubNested struct {
	TrueBool bool              `yaml:"true_bool,omitempty" default:"true"`
	Map      map[string]string `yaml:"map,omitempty" default:"{}"`
}

type Special struct {
	Map     map[string]string `yaml:"map,omitempty" default:"{}"`
	Floats  []float64         `yaml:"floats,omitempty" default:"[]"`
	Strings []string          `yaml:"strings,omitempty" default:"[]"`
}

func LoadConfig(name string) *Config {
	_ = gotenv.Load()

	k := koanf.New(".")

	if name != "" {
		_ = k.Load(file.Provider(name), yaml.Parser())
	}

	_ = k.Load(env2json.Provider(configEnvPrefix, configEnvDelim, func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, configEnvPrefix+configEnvDelim))
	}), json.Parser(), koanf.WithMergeFunc(func(src, dest map[string]interface{}) error {
		return mergo.Merge(&dest, src, mergo.WithSliceDeepCopy)
	}))

	config := &Config{}
	_ = k.UnmarshalWithConf("", config, koanf.UnmarshalConf{
		Tag: "yaml",
		DecoderConfig: &mapstructure.DecoderConfig{
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				defaults.Defaults("yaml", "default"),
				mapstructure.StringToTimeDurationHookFunc(),
				stringtoslice.StringToSlice(),
				stringtomap.StringToMap(),
			),
			Metadata:         nil,
			Result:           config,
			WeaklyTypedInput: true,
		},
	})

	return config
}
