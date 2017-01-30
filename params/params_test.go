package params

import "testing"

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
	if err != nil {
		t.Errorf("New returned error \"%v\"", err)
		return
	}

	// Assert
	if ps.Pattern != "\\.go$" {
		t.Errorf("Expected Pattern to be %q, but got %q\n", "\\.go$", ps.Pattern)
	}
	if ps.Command != "same" {
		t.Errorf("Expected Command to be %q, but got %q\n", "same", ps.Command)
	}
	if ps.Recursive != false {
		t.Errorf("Expected Recursive to be %v, but got %v\n", false, ps.Recursive)
	}
	if ps.IntervalMillis != 500 {
		t.Errorf("Expected Pattern to be %v, but got %v\n", 500, ps.IntervalMillis)
	}
}
