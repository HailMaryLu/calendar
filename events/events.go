package events

import (
	"fmt"
	"time"

	"github.com/araddon/dateparse"
	"github.com/google/uuid"
	"github.com/hailmarylu/calendar/errorLib"
	"github.com/hailmarylu/calendar/reminder"
)

type Event struct {
	ID       string             `json:"id"`
	Title    string             `json:"title"`
	StartAt  time.Time          `json:"start_at"`
	Priority Priority           `json:"priority"`
	Reminder *reminder.Reminder `json:"-"`
}

func getNextID() string {
	return uuid.New().String()
}

func NewEvent(title string, dateStr string, priority string) (*Event, error) {
	if !IsValidTitle(title) {
		return nil, fmt.Errorf("Не удалось создать событие: %w", errorLib.ErrInvalidTitle)
	}
	t, err := dateparse.ParseAny(dateStr)
	if err != nil {
		return nil, fmt.Errorf("Не удалось создать событие: %w", errorLib.ErrInvalidTimeFormat)
	}
	p := Priority(priority)
	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("Не удалось создать событие: %w", err)
	}
	return &Event{
		ID:       getNextID(),
		Title:    title,
		StartAt:  t,
		Priority: p,
		Reminder: nil,
	}, nil
}

func (e *Event) Update(title string, date string, priority string) error {
	if !IsValidTitle(title) {
		return fmt.Errorf("Не удалось обновить событие: %w", errorLib.ErrInvalidTitle)
	}
	t, err := dateparse.ParseAny(date)
	if err != nil {
		return fmt.Errorf("Не удалось обновить событие: %w", errorLib.ErrInvalidTimeFormat)
	}
	p := Priority(priority)
	if err := p.Validate(); err != nil {
		return fmt.Errorf("Не удалось обновить событие: %w", err)
	}
	e.Title = title
	e.StartAt = t
	e.Priority = p
	return nil
}

func (e *Event) AddReminder(message string, at string, notify func(string)) error {
	var err error
	e.Reminder, err = reminder.NewReminder(message, at, notify)
	if err != nil {
		return fmt.Errorf("Не удалось добавить напоминание событию %s: %w", e.Title, err)
	}
	return nil
}

func (e *Event) RemoveReminder() {
	if e.Reminder != nil {
		e.Reminder.Stop()
		e.Reminder = nil
	}
}
