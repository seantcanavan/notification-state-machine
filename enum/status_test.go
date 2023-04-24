package enum

import (
	"github.com/jgroeneveld/trial/assert"
	"testing"
)

func TestStatusStringValues(t *testing.T) {
	expectedScaleTests := []struct {
		name   string
		status Status
		valid  bool
	}{
		{name: "Created", status: "Created", valid: true},
		{name: "Error", status: "Error", valid: true},
		{name: "NA", status: " NA", valid: true},
		{name: "Processing", status: "Processing", valid: true},
		{name: "Queued", status: "Queued", valid: true},
		{name: "Sent", status: "Sent", valid: true},
		{name: "created", status: "created", valid: false},
		{name: "cReAtEd", status: "cReAtEd", valid: false},
		{name: "CREATED", status: "CREATED", valid: false},
		{name: "empty", status: "", valid: false},
	}

	for _, tt := range expectedScaleTests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.status.Valid()
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
func TestStatusMapHasNoNewValues(t *testing.T) {
	assert.Equal(t, len(statusToString), 6)
}
