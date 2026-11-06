package errorLib

import "errors"

var (
	ErrNotExists         = errors.New("ErrNotExists")
	ErrInvalidPriority   = errors.New("ErrInvalidPriority")
	ErrEmptyArchive      = errors.New("ErrEmptyArchive")
	ErrInvalidTitle      = errors.New("ErrInvalidTitle")
	ErrInvalidTimeFormat = errors.New("ErrInvalidTimeFormat")
	ErrReminderNotExists = errors.New("ErrReminderNotExists")
	ErrEmptyMessage      = errors.New("ErrEmptyMessage")
	ErrInvalidTime       = errors.New("ErrInvalidTime")
)

func ReturnCustomerErr(err error) string {
	if errors.Is(err, ErrNotExists) {
		return "События по данному ключу не существует"
	} else if errors.Is(err, ErrInvalidPriority) {
		return "Некорректный приоритет"
	} else if errors.Is(err, ErrEmptyArchive) {
		return "Архив пуст"
	} else if errors.Is(err, ErrInvalidTitle) {
		return "Неверный формат названия"
	} else if errors.Is(err, ErrInvalidTimeFormat) {
		return "Неверный формат даты"
	} else if errors.Is(err, ErrReminderNotExists) {
		return "У события не установлено напоминание"
	} else if errors.Is(err, ErrEmptyMessage) {
		return "Сообщение пустое"
	} else if errors.Is(err, ErrInvalidTime) {
		return "Дата события меньше текущей"
	} else {
		return "Возникла непредвиденная ошибка, обратитесь к администратору -_-"
	}
}
