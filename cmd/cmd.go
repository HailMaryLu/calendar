package cmd

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/google/shlex"
	"github.com/hailmarylu/calendar/calendar"
	"github.com/hailmarylu/calendar/errorLib"
	"github.com/hailmarylu/calendar/logger"
)

type Cmd struct {
	calendar *calendar.Calendar
	iOLog    []LogRecord
	mu       sync.Mutex
}

func NewCmd(c *calendar.Calendar) *Cmd {
	err := c.Load()
	if err != nil {
		fmt.Println(CmdStartErrTitle, err)
		return nil
	}
	return &Cmd{calendar: c}
}

func (c *Cmd) logging(recType string, text string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.iOLog = append(c.iOLog, LogRecord{
		Time: time.Now(),
		Type: recType,
		Text: text,
	})
}

func (c *Cmd) printOutMsg(msgType string, msg string) {
	switch msgType {
	case LogTypeInfo:
		fmt.Println(msg)
	case LogTypeWarn:
		fmt.Println(WarnMsgTitle, msg)
	default:
		fmt.Println(UnknownMsgType)
	}
	c.logging(msgType, msg)
}

func (c *Cmd) printErr(err error) {
	fmt.Println(errorLib.ReturnCustomerErr(err))
	c.logging(LogTypeError, err.Error())
}

func (c *Cmd) showLog() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, record := range c.iOLog {
		fmt.Println(record.Time, " ", record.Type, ": ", record.Text)
	}
}

func (c *Cmd) executor(input string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}
	c.logging(LogTypeInput, input)

	parts, err := shlex.Split(input)
	if err != nil {
		c.printErr(err)
		return
	}
	if len(parts) == 0 {
		return
	}
	cmd := strings.ToLower(parts[0])

	switch cmd {
	case CmdAdd:
		if len(parts) < 4 {
			c.printOutMsg(LogTypeWarn, PrintCmdFormat(cmd))
			return
		}
		title := parts[1]
		date := parts[2]
		priority := parts[3]

		_, err := c.calendar.AddEvent(title, date, priority)
		if err != nil {
			logger.LogError("Ошибка создания события по команде: " + input + "\n" + err.Error())
			c.printErr(err)
		}

	case CmdList:
		list := c.calendar.GetEvents()
		if len(list) == 0 {
			return
		}
		c.printOutMsg(LogTypeInfo, TitleStr)
		for k, e := range list {
			c.printOutMsg(LogTypeInfo, PrintEvent(k, *e))
		}
		c.printOutMsg(LogTypeInfo, DelimLine)

	case CmdRemove:
		if len(parts) < 2 {
			c.printOutMsg(LogTypeWarn, PrintCmdFormat(cmd))
			return
		}
		id := parts[1]
		err := c.calendar.DeleteEvent(id)
		if err != nil {
			logger.LogError("Ошибка удаления события по команде: " + input + "\n" + err.Error())
			c.printErr(err)
		}

	case CmdAddReminder:
		if len(parts) < 4 {
			c.printOutMsg(LogTypeWarn, PrintCmdFormat(cmd))
			return
		}
		id := parts[1]
		msg := parts[2]
		timeAt := parts[3]
		err := c.calendar.SetEventReminder(id, msg, timeAt)
		if err != nil {
			logger.LogError("Ошибка установки напоминания по команде: " + input + "\n" + err.Error())
			c.printErr(err)
		}
	case CmdDelReminder:
		if len(parts) < 2 {
			c.printOutMsg(LogTypeWarn, PrintCmdFormat(cmd))
			return
		}
		id := parts[1]
		err := c.calendar.CancelEventReminder(id)
		if err != nil {
			logger.LogError("Ошибка удаления напоминания по команде: " + input + "\n" + err.Error())
			c.printErr(err)
		}
	case CmdLog:
		c.showLog()

	case CmdHelp:
		c.printOutMsg(LogTypeInfo, PrintHelp())

	case CmdExit:
		err := c.calendar.Exit()
		if err != nil {
			logger.LogError("Ошибка завершения работы по команде: " + input + "\n" + err.Error())
			c.printErr(err)
		}
		logger.LogInfo("Выход из программы")
		err = logger.CloseFile()
		if err != nil {
			logger.LogError("Ошибка закрытия файла журналирования\n" + err.Error())
			c.printErr(err)
		}
		os.Exit(0)

	default:
		c.printOutMsg(LogTypeWarn, fmt.Sprint("Неизвестная команда: ", cmd, ". Введите ", CmdHelp, " для списка команд"))
	}

}

func (c *Cmd) completer(d prompt.Document) []prompt.Suggest {
	suggestions := []prompt.Suggest{
		{Text: CmdAdd, Description: CmdDescAdd},
		{Text: CmdList, Description: CmdDescList},
		{Text: CmdRemove, Description: CmdDescRemove},
		{Text: CmdAddReminder, Description: CmdDescAddRem},
		{Text: CmdDelReminder, Description: CmdDescDelRem},
		{Text: CmdLog, Description: CmdDescLog},
		{Text: CmdHelp, Description: CmdDescHelp},
		{Text: CmdExit, Description: CmdDescExit},
	}
	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}

func (c *Cmd) Run() {
	c.logging(LogTypeInfo, "App starting")
	logger.LogInfo("Запуск консоли")
	go func() {
		for msg := range c.calendar.Notification {
			c.printOutMsg(LogTypeInfo, fmt.Sprintln(NotifMsgTitle, msg))
		}
	}()
	p := prompt.New(c.executor, c.completer, prompt.OptionPrefix("> "))
	p.Run()
}
