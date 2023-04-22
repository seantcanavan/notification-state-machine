package queue

import (
	"github.com/jgroeneveld/trial/assert"
	"testing"
)

func TestMessageTypeStringValues(t *testing.T) {
	messageTypeTests := []struct {
		name        string
		messageType MessageType
		valid       bool
	}{
		{name: "CreateReq", messageType: "CreateReq", valid: true},
		{name: "CrEaTeReQ", messageType: "CrEaTeReQ", valid: false},
		{name: "CREATEREQ", messageType: "CREATEREQ", valid: false},
		{name: "empty", messageType: "", valid: false},
	}

	for _, tt := range messageTypeTests {
		// using tt.name from the case to use it as the `t.Run` test name
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.messageType.Valid()
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
func TestMessageTypeMapHasNoNewValues(t *testing.T) {
	assert.Equal(t, len(messageTypeToString), 1)
}
