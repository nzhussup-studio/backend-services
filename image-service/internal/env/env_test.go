package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	t.Setenv("TEST_STRING", "hello")

	assert.Equal(t, "hello", GetString("TEST_STRING", "default"))
	assert.Equal(t, "fallback", GetString("MISSING_STRING", "fallback"))
}

func TestGetInt(t *testing.T) {
	t.Setenv("TEST_INT", "42")

	assert.Equal(t, 42, GetInt("TEST_INT", 99))
	assert.Equal(t, 99, GetInt("MISSING_INT", 99))

	t.Setenv("TEST_INT_INVALID", "notanint")
	assert.Equal(t, 88, GetInt("TEST_INT_INVALID", 88))
}

func TestGetBool(t *testing.T) {
	t.Setenv("TEST_BOOL_TRUE", "true")
	t.Setenv("TEST_BOOL_FALSE", "false")

	assert.True(t, GetBool("TEST_BOOL_TRUE", false))
	assert.False(t, GetBool("TEST_BOOL_FALSE", true))
	assert.True(t, GetBool("MISSING_BOOL", true))

	t.Setenv("TEST_BOOL_INVALID", "notabool")
	assert.False(t, GetBool("TEST_BOOL_INVALID", false))
}
