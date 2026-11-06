package storage

import (
	"archive/zip"
	"fmt"
	"io"
	"os"

	"github.com/hailmarylu/calendar/errorLib"
)

type ZipStorage struct {
	*Storage
}

func NewZipStorage(filename string) *ZipStorage {
	return &ZipStorage{
		&Storage{filename: filename},
	}
}

func (z *ZipStorage) Save(data []byte) error {
	f, err := os.Create(z.GetFileName())
	if err != nil {
		return fmt.Errorf("Не удалось создать файл  для архива: %w", err)
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	defer zw.Close()

	w, err := zw.Create("data")
	if err != nil {
		return fmt.Errorf("Не удалось создать архив: %w", err)
	}
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("Не удалось записать данные в архив: %w", err)
	}
	return nil
}

func (z *ZipStorage) Load() ([]byte, error) {
	r, err := zip.OpenReader(z.GetFileName())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("Не удалось открыть архив: %w", err)
	}

	if len(r.File) == 0 {
		return nil, errorLib.ErrEmptyArchive
	}

	file := r.File[0]
	rc, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("Не удалось открыть файл из архива: %w", err)
	}

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("Не удалось прочитать файл из архива: %w", err)
	}

	err = rc.Close()
	if err != nil {
		return nil, fmt.Errorf("Не удалось закрыть файл из архива: %w", err)
	}
	err = r.Close()
	if err != nil {
		return nil, fmt.Errorf("Не удалось закрыть архив: %w", err)
	}
	return data, nil
}
