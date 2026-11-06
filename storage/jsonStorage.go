package storage

import (
	"fmt"
	"os"
)

type JsonStorage struct {
	*Storage
}

func NewStorage(filename string) *JsonStorage {
	return &JsonStorage{
		&Storage{filename: filename},
	}
}

func (s *JsonStorage) Save(data []byte) error {
	err := os.WriteFile(s.GetFileName(), data, 0644)
	if err != nil {
		return fmt.Errorf("Не удалось записать данные в JSON: %w", err)
	}
	return nil
}

func (s *JsonStorage) Load() ([]byte, error) {
	res, err := os.ReadFile(s.GetFileName())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("Не удалось прочитать JSON: %w", err)
	}
	return res, nil
}
