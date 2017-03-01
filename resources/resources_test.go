package resources

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testResourceFile = "./example.rc"

func TestNew(t *testing.T) {
	// Act
	rcs, err := New([]string{"unexisting-file", testResourceFile})

	// Assert
	assert.Nil(t, err, fmt.Sprintf("New returned error \"%v\"", err))
	assert.Equal(t, len(rcs.Map), 5, fmt.Sprintf("Expected 5 resources, but got %v", len(rcs.Map)))
	assertSetting(t, rcs.Map, "what-file", ".examplerc")
	assertSetting(t, rcs.Map, "pattern", `\.go$`)
	assertSetting(t, rcs.Map, "command", "go test ./...")
}

func assertSetting(t *testing.T, m map[string]string, k, e string) {
	v, ok := m[k]
	if !assert.True(t, ok, fmt.Sprintf("Resource %q not found", k)) {
		return
	}
	assert.Equal(t, v, e, fmt.Sprintf("Resource %q => %q, but expected %q", k, v, e))
}
