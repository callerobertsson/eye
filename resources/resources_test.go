package resources

import "testing"

const testResourceFile = "./example.rc"

func TestNew(t *testing.T) {
	rcs, err := New([]string{"unexisting-file", testResourceFile})

	if err != nil {
		t.Errorf("New returned error \"%v\"", err)
		return
	}

	if len(rcs.Map) != 5 {
		t.Errorf("Expected 5 resources, but got %v", len(rcs.Map))
	}

	validateSetting(t, rcs.Map, "what-file", ".examplerc")
	validateSetting(t, rcs.Map, "pattern", `\.go$`)
	validateSetting(t, rcs.Map, "command", "go test ./...")

}

func validateSetting(t *testing.T, m map[string]string, k, e string) {
	v, ok := m[k]
	if !ok {
		t.Errorf("Resource %q was not found", k)
	}
	if v != e {
		t.Errorf("Resource %q => %q, but expected %q", k, v, e)
	}
}
