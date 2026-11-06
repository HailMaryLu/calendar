package events

import "testing"

func TestValidate(t *testing.T) {
	p := Priority("low")
	err := p.Validate()
	if err != nil {
		t.Errorf("Expected no error for priority '%v', err: %v", p, err.Error())
	}

	p = Priority("critical")
	err = p.Validate()
	if err == nil {
		t.Errorf("Expected error for priority '%v'", p)
	}
}
