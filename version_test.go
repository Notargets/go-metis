package metis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	v := Version()
	assert.NotEmpty(t, v)
	assert.Equal(t, "5.2.1", v)
}

func TestGoMetisVersion(t *testing.T) {
	v := GoMetisVersion()
	assert.NotEmpty(t, v)
	// Will be "dev" in development, or the actual tag when exported
	t.Logf("go-metis version: %s", v)
}
