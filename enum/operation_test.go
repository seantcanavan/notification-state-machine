package enum

import (
	"github.com/jgroeneveld/trial/assert"
	"testing"
)

func TestOperationStringValues(t *testing.T) {
	expectedScaleTests := []struct {
		name      string
		operation Operation
		valid     bool
	}{
		{name: "Create", operation: "Create", valid: true},
		{name: "Delete", operation: "Delete", valid: true},
		{name: "Read", operation: "Read", valid: true},
		{name: "Update", operation: "Update", valid: true},
		{name: "Permanent", operation: "Permanent", valid: false},
		{name: "create", operation: "create", valid: false},
		{name: "cReAtE", operation: "cReAtE", valid: false},
		{name: "CREATE", operation: "CREATE", valid: false},
		{name: "empty", operation: "", valid: false},
	}

	for _, tt := range expectedScaleTests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.operation.Valid()
			if valid != tt.valid {
				t.Errorf("got %t want %t", valid, tt.valid)
			}
		})
	}
}

// If this test fails then you need to expand the above
// tests because you added new values to the map without
// testing them in the above test! This guarantees consistency
// between the map that reverses the values and the values which
// are also checked against the map for validity.
func TestOperationMapHasNoNewValues(t *testing.T) {
	assert.Equal(t, len(operationToString), 4)
}
