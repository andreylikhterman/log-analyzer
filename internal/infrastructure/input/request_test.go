package input

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckFlags_Valid(t *testing.T) {
	requestTemplate := RequestTemplate{
		RequeredFlags: []string{"path"},
		OptionalFlags: []string{"from", "to", "format", "filter-field", "filter-value"},
	}

	parts := []string{"--path", "/some/path", "--from", "2022-01-01"}
	require.NoError(t, checkFlags(requestTemplate, parts))
}

func TestCheckFlags_Invalid(t *testing.T) {
	requestTemplate := RequestTemplate{
		RequeredFlags: []string{"path"},
		OptionalFlags: []string{"from", "to", "format", "filter-field", "filter-value"},
	}

	parts := []string{"--unknown", "value"}
	assert.Error(t, checkFlags(requestTemplate, parts))
}

func TestCheckCountFlags_Valid(t *testing.T) {
	requestTemplate := RequestTemplate{
		RequeredFlags: []string{"path"},
		OptionalFlags: []string{"from", "to", "format", "filter-field", "filter-value"},
	}

	parts := []string{"--path", "/some/path", "--from", "2022-01-01"}
	require.NoError(t, checkCountFlags(requestTemplate, parts))
}

func TestCheckCountFlags_Invalid(t *testing.T) {
	requestTemplate := RequestTemplate{
		RequeredFlags: []string{"path"},
		OptionalFlags: []string{"from", "to", "format", "filter-field", "filter-value"},
	}

	parts := []string{"--path", "/some/path", "--path", "/another/path"}
	assert.Error(t, checkCountFlags(requestTemplate, parts))
}

func TestGetFlags_Valid(t *testing.T) {
	requestTemplate := RequestTemplate{
		RequeredFlags: []string{"path"},
		OptionalFlags: []string{"from", "to", "format", "filter-field", "filter-value"},
	}

	parts := []string{"--path", "/some/path", "--from", "2022-01-01"}
	flags := getFlags(requestTemplate, parts)

	assert.Equal(t, "/some/path", flags["path"])
	assert.Equal(t, "2022-01-01", flags["from"])
	assert.Empty(t, flags["to"])
}
