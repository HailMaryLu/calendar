package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	file   *os.File
	logger *log.Logger
	mu     sync.Mutex
)

const (
	dir        = "logs"
	dateFormat = "20060102_150405.000"
	fileFormat = ".log"
)

func InitializeLogger() {
	mu.Lock()
	defer mu.Unlock()

	if logger != nil {
		return
	}

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Println("Ошибка создания директории логов:", err.Error())
		return
	}
	filename := time.Now().Format(dateFormat) + fileFormat
	finalPath := filepath.Join(dir, filename)
	file, err = os.OpenFile(finalPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Ошибка создания файла для журналирования:", err.Error())
		return
	}
	logger = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func LogInfo(txt string) {
	mu.Lock()
	defer mu.Unlock()
	if logger != nil {
		logger.SetPrefix("INFO: ")
		err := logger.Output(2, txt)
		if err != nil {
			fmt.Println("Ошибка при записи info сообщения:", err.Error())
		}
	}
}

func LogError(txt string) {
	mu.Lock()
	defer mu.Unlock()
	if logger != nil {
		logger.SetPrefix("ERROR: ")
		err := logger.Output(2, txt)
		if err != nil {
			fmt.Println("Ошибка при записи error сообщения:", err.Error())
		}
	}
}

func CloseFile() error {
	mu.Lock()
	defer mu.Unlock()

	if file != nil {
		return file.Close()
	}
	return nil
}
