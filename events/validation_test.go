package events

import (
	"testing"
)

func TestIsValidTitle(t *testing.T) {
	title := "some correct title"
	result := IsValidTitle(title)
	if result != true {
		t.Errorf("Expected true for title %s", title)
	}

	title = "!ncorrect title"
	result = IsValidTitle(title)
	if result != false {
		t.Errorf("Expected false for title %s", title)
	}
}
