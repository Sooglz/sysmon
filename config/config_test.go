package config

import "testing"

func TestDefaultValid(t *testing.T) {
	c := Default()
	if err := c.Validate(); err != nil {
		t.Fatalf("default config should be valid: %v", err)
	}
}

func TestValidateErrors(t *testing.T) {
	cases := []struct {
		name string
		mod  func(*Config)
	}{
		{"interval zero", func(c *Config) { c.Interval = 0 }},
		{"empty listen", func(c *Config) { c.ListenAddr = "" }},
		{"empty disk", func(c *Config) { c.DiskPath = "" }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := Default()
			tc.mod(&c)
			if err := c.Validate(); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
