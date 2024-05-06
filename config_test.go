package bootstrap

import "testing"

func TestRegisterConfig(t *testing.T) {
	var cfg struct {
		Test string
	}
	RegisterConfig(&cfg)
}
