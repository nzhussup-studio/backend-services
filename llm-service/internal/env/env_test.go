package env

import (
	"os"
	"testing"
)

func TestGetString(t *testing.T) {
	tests := map[string]struct {
		key          string
		envValue     string
		defaultValue string
		want         string
	}{
		"env var set": {
			key: "TEST_STRING_SET", envValue: "hello", defaultValue: "default", want: "hello",
		},
		"env var empty": {
			key: "TEST_STRING_EMPTY", envValue: "", defaultValue: "default", want: "default",
		},
		"env var unset": {
			key: "TEST_STRING_UNSET", defaultValue: "default", want: "default",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			got := GetString(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("GetString(%q, %q) = %q; want %q", tt.key, tt.defaultValue, got, tt.want)
			}
		})
	}
}

func TestGetInt(t *testing.T) {
	tests := map[string]struct {
		key          string
		envValue     string
		defaultValue int
		want         int
	}{
		"valid int": {
			key: "TEST_INT", envValue: "42", defaultValue: 0, want: 42,
		},
		"invalid int": {
			key: "TEST_INT_INVALID", envValue: "abc", defaultValue: 10, want: 10,
		},
		"empty value": {
			key: "TEST_INT_EMPTY", envValue: "", defaultValue: 5, want: 5,
		},
		"unset": {
			key: "TEST_INT_UNSET", defaultValue: 100, want: 100,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			got := GetInt(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("GetInt(%q, %d) = %d; want %d", tt.key, tt.defaultValue, got, tt.want)
			}
		})
	}
}

func TestGetBool(t *testing.T) {
	tests := map[string]struct {
		key          string
		envValue     string
		defaultValue bool
		want         bool
	}{
		"true value": {
			key: "TEST_BOOL_TRUE", envValue: "true", defaultValue: false, want: true,
		},
		"false value": {
			key: "TEST_BOOL_FALSE", envValue: "false", defaultValue: true, want: false,
		},
		"invalid bool": {
			key: "TEST_BOOL_INVALID", envValue: "maybe", defaultValue: true, want: true,
		},
		"empty value": {
			key: "TEST_BOOL_EMPTY", envValue: "", defaultValue: false, want: false,
		},
		"unset": {
			key: "TEST_BOOL_UNSET", defaultValue: true, want: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			got := GetBool(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("GetBool(%q, %v) = %v; want %v", tt.key, tt.defaultValue, got, tt.want)
			}
		})
	}
}
