package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_LoadConfig(t *testing.T) {
	t.Setenv("APP__STRING", "hello")
	t.Setenv("APP__STRUCT__TRUE_BOOL", "false")
	t.Setenv("APP__STRUCT__FALSE_BOOL", "true")
	t.Setenv("APP__SLICE__0__TRUE_BOOL", "false")
	t.Setenv("APP__SLICE__1__FALSE_BOOL", "true")

	want := &Config{
		String: "hello",
		Bool:   false,
		Struct: Nested{
			TrueBool:  false,
			FalseBool: true,
		},
		Slice: []Nested{
			{
				TrueBool:  false,
				FalseBool: false,
			},
			{
				TrueBool:  false,
				FalseBool: true,
			},
		},
	}
	cfg := LoadConfig()

	assert.Equal(t, want, cfg, "LoadConfig()")
}
