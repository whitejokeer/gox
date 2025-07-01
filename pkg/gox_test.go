package gox

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	version := Version()
	assert.NotEmpty(t, version)
}

func TestFramework_New(t *testing.T) {
	framework, err := New("/tmp/output")

	require.NoError(t, err)
	assert.NotNil(t, framework)

	// Clean up
	err = framework.Close()
	assert.NoError(t, err)
}
