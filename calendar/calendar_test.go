package calendar

import (
	"testing"

	"github.com/hailmarylu/calendar/storage"
)

func TestDeleteEvent(t *testing.T) {
	zs := storage.NewZipStorage("fileName")
	cal := NewCalendar(zs)
	id := "67"
	err := cal.DeleteEvent(id)
	if err == nil {
		t.Errorf("Expected error for event id '%v': not exists", id)
	}
}

func TestAddEvent(t *testing.T) {
	zs := storage.NewZipStorage("fileName")
	cal := NewCalendar(zs)
	_, err := cal.AddEvent("title", "2026-02-29 10:00:00", "low")
	if err == nil {
		t.Error("Expected error for incorrect date")
	}

	e, _ := cal.AddEvent("title", "2026-02-28 10:00:00", "low")
	if e == nil {
		t.Error("Expected no error on adding event")
	}
}
