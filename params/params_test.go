package params

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testResourceFile = "./example.rc"

func TestAddParamsFromResourceFile(t *testing.T) {
	// Arrange
	ps := Params{
		Pattern:        "change", // should change to \.go$
		Command:        "same",   // should not change
		Recursive:      true,     // should change to false
		IntervalMillis: 1000,     // should change to 500
	}

	// Act
	err := ps.AddParamsFromResourceFile([]string{"unexisting-file", testResourceFile})
	assert.Nil(t, err, fmt.Sprintf("New returned error \"%v\"", err))

	// Assert
	assert.Equal(t, ps.Pattern, "\\.go$", "Wrong pattern")
	assert.Equal(t, ps.Command, "same", "Command mismatch")
	assert.False(t, ps.Recursive, "Recursive flag error")
	assert.Equal(t, ps.IntervalMillis, 500, "Wrong IntervalMillis")
}
