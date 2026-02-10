package age

import (
	"filippo.io/age/plugin"
	"testing"
)

var testOnlyAgePassword string

func printf(format string, v ...interface{}) {
	log.Printf("age: "+format, v...)
}

func warningf(format string, v ...interface{}) {
	log.Printf("age: warning: "+format, v...)
}

var pluginTerminalUIImpl = plugin.NewTerminalUI(printf, warningf)

// We cannot use plugin.NewTerminalUI() directly because we want to be able to
// inject specific return values for RequestValue during testing.
var pluginTerminalUI = &plugin.ClientUI{
	DisplayMessage: func(name, message string) error {
		return pluginTerminalUIImpl.DisplayMessage(name, message)
	},
	RequestValue: func(name, message string, isSecret bool) (s string, err error) {
		if testing.Testing() && testOnlyAgePassword != "" {
			return testOnlyAgePassword, nil
		}
		return pluginTerminalUIImpl.RequestValue(name, message, isSecret);
	},
	Confirm: func(name, message, yes, no string) (choseYes bool, err error) {
		return pluginTerminalUIImpl.Confirm(name, message, yes, no)
	},
	WaitTimer: func(name string) {
		pluginTerminalUIImpl.WaitTimer(name)
	},
}
