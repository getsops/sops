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

func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	bytes, err := f.TextFormatter.Format(entry)
	name := color.New(color.Bold).Sprintf("[%s]", f.LoggerName)
	return []byte(fmt.Sprintf("%s\t %s", name, bytes)), err
}

func NewLogger(name string) *logrus.Logger {
	log := logrus.New()
	log.Formatter = &TextFormatter{
		LoggerName: name,
	}
	Loggers[name] = log
	return log
}

var Loggers map[string]*logrus.Logger
