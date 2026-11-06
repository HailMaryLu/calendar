package main

import (
	"os"

	"github.com/hailmarylu/calendar/calendar"
	"github.com/hailmarylu/calendar/cmd"
	"github.com/hailmarylu/calendar/logger"
	"github.com/hailmarylu/calendar/storage"
)

const fileName = "calendar.zip"

func main() {
	logger.InitializeLogger()
	//s := storage.NewStorage("calendar.json")
	zs := storage.NewZipStorage(fileName)
	cal := calendar.NewCalendar(zs)

	cli := cmd.NewCmd(cal)
	if cli == nil {
		logger.LogError("Ошибка запуска консоли")
		os.Exit(1)
	}
	cli.Run()
}
