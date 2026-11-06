package calendar

import (
	"encoding/json"
	"fmt"

	"github.com/hailmarylu/calendar/errorLib"
	"github.com/hailmarylu/calendar/events"
	"github.com/hailmarylu/calendar/logger"
	"github.com/hailmarylu/calendar/storage"
)

type Calendar struct {
	CalendarEvents map[string]*events.Event `json:"events"`
	storage        storage.Store
	Notification   chan string `json:"-"`
}

func NewCalendar(s storage.Store) *Calendar {
	logger.LogInfo("Создание календаря")
	return &Calendar{
		CalendarEvents: make(map[string]*events.Event),
		storage:        s,
		Notification:   make(chan string),
	}
}

func (c *Calendar) getEvent(key string) *events.Event {
	return c.CalendarEvents[key]
}

func (c *Calendar) AddEvent(title string, date string, priority string) (*events.Event, error) {
	e, err := events.NewEvent(title, date, priority)
	if err != nil {
		return nil, err
	}
	c.CalendarEvents[e.ID] = e
	return e, nil
}

func (c *Calendar) EditEvent(key string, newTitle string, newDate string, newPriority string) error {
	e := c.getEvent(key)
	if e == nil {
		return fmt.Errorf("Не удалось отредактировать событие: %w", errorLib.ErrNotExists)
	}
	return e.Update(newTitle, newDate, newPriority)
}

func (c *Calendar) DeleteEvent(key string) error {
	e := c.getEvent(key)
	if e == nil {
		return fmt.Errorf("Не удалось удалить событие: %w", errorLib.ErrNotExists)
	}
	delete(c.CalendarEvents, key)
	return nil
}

func (c *Calendar) GetEvents() map[string]*events.Event {
	return c.CalendarEvents
}

func (c *Calendar) Save() error {
	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("Не удалось упаковать календарь в файл: %w", err)
	}
	return c.storage.Save(data)
}

func (c *Calendar) Load() error {
	data, err := c.storage.Load()
	if err != nil {
		return fmt.Errorf("Не удалось достать данные из хранилища: %w", err)
	}
	if data != nil {
		err = json.Unmarshal(data, &c)
	}
	if err != nil {
		return fmt.Errorf("Не удалось разобрать файл с календарём: %w", err)
	}
	return nil
}

func (c *Calendar) Notify(msg string) {
	logger.LogInfo("Отправка сообщения в канал нотификаций")
	c.Notification <- msg
}

func (c *Calendar) SetEventReminder(id string, msg string, at string) error {
	e := c.getEvent(id)
	if e == nil {
		return fmt.Errorf("Не удалось получить событие для добавления напоминания: %w", errorLib.ErrNotExists)
	}
	err := e.AddReminder(msg, at, c.Notify)
	if err != nil {
		return err
	}
	return nil
}

func (c *Calendar) CancelEventReminder(id string) error {
	e := c.getEvent(id)
	if e == nil {
		return fmt.Errorf("Не удалось удалить событие: %w", errorLib.ErrNotExists)
	}
	e.RemoveReminder()
	return nil
}

func (c *Calendar) Exit() error {
	logger.LogInfo("Закрытие календаря")
	for _, ev := range c.CalendarEvents {
		ev.RemoveReminder()
	}
	err := c.Save()
	if err != nil {
		return fmt.Errorf("Не удалось сохранить календарь при завершении работы: %w", err)
	}
	close(c.Notification)
	return nil
}
