package reminder

import (
	"fmt"
	"time"

	"github.com/araddon/dateparse"
	"github.com/hailmarylu/calendar/errorLib"
)

type Reminder struct {
	Message string
	At      time.Time
	Sent    bool
	Timer   *time.Timer
	Notify  func(string)
}

func isMessageEmpty(m string) bool {
	return m == ""
}

func isDateValid(t time.Time) bool {
	return time.Now().Before(t)
}

func NewReminder(message string, at string, notify func(string)) (*Reminder, error) {
	if isMessageEmpty(message) {
		return nil, fmt.Errorf("Не удалось создать напоминание: %w", errorLib.ErrEmptyMessage)
	}
	t, err := dateparse.ParseAny(at)
	if err != nil {
		return nil, fmt.Errorf("Не удалось создать напоминание: %w", errorLib.ErrInvalidTimeFormat)
	}
	if !isDateValid(t) {
		return nil, fmt.Errorf("Не удалось создать напоминание: %w", errorLib.ErrInvalidTime)
	}

	r := &Reminder{
		Message: message,
		At:      t,
		Sent:    false,
		Notify:  notify,
	}

	dur := time.Until(t)
	r.Timer = time.AfterFunc(dur, r.Start)
	return r, nil
}

func (r *Reminder) Send() {
	if r.Sent {
		return
	}
	if r.Notify != nil {
		r.Notify(r.Message)
		r.Sent = true
	}
}

func (r *Reminder) Start() {
	r.Send()
}

func (r *Reminder) Stop() {
	if r.Timer != nil {
		r.Timer.Stop()
	}
}
