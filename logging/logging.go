package logging

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

func init() {
	Loggers = make(map[string]*logrus.Logger)
}

// TextFormatter extends the standard logrus TextFormatter and adds a field to specify the logger's name
type TextFormatter struct {
	LoggerName string
	logrus.TextFormatter
}

// Format formats a log entry onto bytes
func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	bytes, err := f.TextFormatter.Format(entry)
	name := color.New(color.Bold).Sprintf("[%s]", f.LoggerName)
	return []byte(fmt.Sprintf("%s\t %s", name, bytes)), err
}

// NewLogger is the constructor for a new Logger object with the given name
func NewLogger(name string) *logrus.Logger {
	log := logrus.New()
	log.SetLevel(logrus.WarnLevel)
	log.Formatter = &TextFormatter{
		LoggerName: name,
	}
	Loggers[name] = log
	return log
}

// SetLevel sets the given level for all current Loggers
func SetLevel(level logrus.Level) {
	for k := range Loggers {
		Loggers[k].SetLevel(level)
	}
}

// Loggers is the runtime map of logger name to logger object
var Loggers map[string]*logrus.Logger
