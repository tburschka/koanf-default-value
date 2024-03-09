package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_LoadConfig_NoConfig(t *testing.T) {
	want := &Config{
		String: "string",
		Bool:   false,
		Struct: Nested{
			TrueBool:  true,
			FalseBool: false,
			SubBool: SubNested{
				TrueBool: true,
				Map:      map[string]string{},
			},
		},
		Slice: []Nested{},
		Special: Special{
			Map:     map[string]string{},
			Floats:  []float64{},
			Strings: []string{},
		},
	}
	cfg := LoadConfig("")

	assert.Equal(t, want, cfg, "LoadConfig()")
}

func Test_LoadConfig_Env(t *testing.T) {
	t.Setenv("APP__STRUCT__TRUE_BOOL", "false")
	t.Setenv("APP__STRUCT__FALSE_BOOL", "true")
	t.Setenv("APP__SLICE__0__TRUE_BOOL", "false")
	t.Setenv("APP__SLICE__0__SUB_BOOL__TRUE_BOOL", "false")
	t.Setenv("APP__SLICE__1__TRUE_BOOL", "false")
	t.Setenv("APP__SLICE__1__FALSE_BOOL", "true")
	t.Setenv("APP__SPECIAL__FLOATS", "[ .1, .5 ]")
	t.Setenv("APP__SPECIAL__STRINGS", "[ 'abc', 'xyz' ]")
	t.Setenv("APP__SPECIAL__MAP", "{'a':'x','b':'y','c':'z'}")

	want := &Config{
		String: "string",
		Bool:   false,
		Struct: Nested{
			TrueBool:  false,
			FalseBool: true,
			SubBool: SubNested{
				TrueBool: true,
				Map:      map[string]string{},
			},
		},
		Slice: []Nested{
			{
				TrueBool:  false,
				FalseBool: false,
				SubBool: SubNested{
					TrueBool: false,
					Map:      map[string]string{},
				},
			},
			{
				TrueBool:  false,
				FalseBool: true,
				SubBool: SubNested{
					TrueBool: true,
					Map:      map[string]string{},
				},
			},
		},
		Special: Special{
			Map: map[string]string{
				"a": "x",
				"b": "y",
				"c": "z",
			},
			Floats:  []float64{.1, .5},
			Strings: []string{"abc", "xyz"},
		},
	}
	cfg := LoadConfig("")

	assert.Equal(t, want, cfg, "LoadConfig()")
}

func Test_LoadConfig_File(t *testing.T) {
	want := &Config{
		String: "string",
		Bool:   false,
		Struct: Nested{
			TrueBool:  false,
			FalseBool: true,
			SubBool: SubNested{
				TrueBool: true,
				Map:      map[string]string{},
			},
		},
		Slice: []Nested{
			{
				TrueBool:  false,
				FalseBool: false,
				SubBool: SubNested{
					TrueBool: false,
					Map: map[string]string{
						"a": "x",
						"b": "y",
						"c": "z",
					},
				},
			},
			{
				TrueBool:  false,
				FalseBool: true,
				SubBool: SubNested{
					TrueBool: true,
					Map:      map[string]string{},
				},
			},
		},
		Special: Special{
			Map: map[string]string{
				"a": "x",
				"b": "y",
				"c": "z",
			},
			Floats:  []float64{.1, .5},
			Strings: []string{"abc", "xyz"},
		},
	}
	cfg := LoadConfig("config.yml")

	assert.Equal(t, want, cfg, "LoadConfig()")
}

func Test_LoadConfig_Mixed(t *testing.T) {
	t.Setenv("APP__STRING", "strong")
	t.Setenv("APP__STRUCT__FALSE_BOOL", "false")
	t.Setenv("APP__SLICE__0__TRUE_BOOL", "true")
	t.Setenv("APP__SLICE__0__SUB_BOOL__FALSE_BOOL", "true")
	t.Setenv("APP__SLICE__1__TRUE_BOOL", "false")
	t.Setenv("APP__SLICE__1__FALSE_BOOL", "true")
	t.Setenv("APP__SPECIAL__FLOATS", "[ .1, .5 ]")
	t.Setenv("APP__SPECIAL__STRINGS", "[ 'abc', 'xyz' ]")
	t.Setenv("APP__SPECIAL__MAP", "{'x':'a','y':'b','z':'c'}")

	want := &Config{
		String: "strong",
		Bool:   false,
		Struct: Nested{
			TrueBool:  false,
			FalseBool: false,
			SubBool: SubNested{
				TrueBool: true,
				Map:      map[string]string{},
			},
		},
		Slice: []Nested{
			{
				TrueBool:  true,
				FalseBool: false,
				SubBool: SubNested{
					TrueBool: false,
					Map: map[string]string{
						"a": "x",
						"b": "y",
						"c": "z",
					},
				},
			},
			{
				TrueBool:  false,
				FalseBool: true,
				SubBool: SubNested{
					TrueBool: true,
					Map:      map[string]string{},
				},
			},
		},
		Special: Special{
			Map: map[string]string{
				"x": "a",
				"y": "b",
				"z": "c",
			},
			Floats:  []float64{.1, .5},
			Strings: []string{"abc", "xyz"},
		},
	}
	cfg := LoadConfig("config.yml")

	assert.Equal(t, want, cfg, "LoadConfig()")
}
