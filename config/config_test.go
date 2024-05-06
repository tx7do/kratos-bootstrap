package config

import "testing"

func TestRegisterConfig(t *testing.T) {
	var cfg struct {
		Test string
	}
	RegisterConfig(&cfg)
}
