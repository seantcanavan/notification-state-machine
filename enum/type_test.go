package enum

import (
	"github.com/jgroeneveld/trial/assert"
	"testing"
)

func TestTypeStringValues(t *testing.T) {
	expectedScaleTests := []struct {
		name  string
		Type  Type
		valid bool
	}{
		{name: "Email", Type: "Email", valid: true},
		{name: "SMS", Type: "SMS", valid: true},
		{name: "Snail", Type: "Snail", valid: true},
		{name: "email", Type: "email", valid: false},
		{name: "eMaIl", Type: "eMaIl", valid: false},
		{name: "EMAIL", Type: "EMAIL", valid: false},
		{name: "empty", Type: "", valid: false},
	}

	for _, tt := range expectedScaleTests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.Type.Valid()
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
func TestTypeMapHasNoNewValues(t *testing.T) {
	assert.Equal(t, len(statusToString), 5)
}
