package cmd

import (
	"fmt"

	"github.com/hailmarylu/calendar/events"
)

func PrintEvent(key string, e events.Event) string {
	return fmt.Sprintf("ID:%s Event:%s Start at:%v", key, e.Title, e.StartAt)
}

func PrintCmdFormat(cmd string) string {
	var s string
	switch cmd {
	case CmdAdd:
		s = fmt.Sprintf("Формат передачи команды: %s \"название события\" \"дата и время\" \"приоритет\"", cmd)
	case CmdRemove:
		s = fmt.Sprintf("Формат передачи команды: %s \"идентификатор события\"", cmd)
	case CmdAddReminder:
		s = fmt.Sprintf("Формат передачи команды: %s \"идентификатор события\"  \"текст напоминания\" \"дата и время срабатывания\"", cmd)
	case CmdDelReminder:
		s = fmt.Sprintf("Формат передачи команды: %s \"идентификатор события\"", cmd)
	}
	return s
}

func PrintHelp() string {
	s := fmt.Sprintln(CmdAdd, "-", CmdDescAdd) +
		PrintCmdFormat(CmdAdd) + "\n" +
		fmt.Sprintln(CmdList, "-", CmdDescList) +
		fmt.Sprintln(CmdRemove, "-", CmdDescRemove) +
		PrintCmdFormat(CmdRemove) + "\n" +
		fmt.Sprintln(CmdAddReminder, "-", CmdDescAddRem) +
		PrintCmdFormat(CmdAddReminder) + "\n" +
		fmt.Sprintln(CmdDelReminder, "-", CmdDescDelRem) +
		PrintCmdFormat(CmdDelReminder) + "\n" +
		fmt.Sprintln(CmdLog, "-", CmdDescLog) +
		fmt.Sprintln(CmdHelp, "-", CmdDescHelp) +
		fmt.Sprint(CmdExit, " - ", CmdDescExit)
	return s
}
