package gobrake

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

var defaultContext map[string]interface{}

func getDefaultContext() map[string]interface{} {
	if defaultContext != nil {
		return defaultContext
	}

	defaultContext = map[string]interface{}{
		"notifier": map[string]interface{}{
			"name":    "gobrake",
			"version": "2.0.4",
			"url":     "https://github.com/airbrake/gobrake",
		},

		"language":     runtime.Version(),
		"os":           runtime.GOOS,
		"architecture": runtime.GOARCH,
	}
	if s, err := os.Hostname(); err == nil {
		defaultContext["hostname"] = s
	}
	if s := os.Getenv("GOPATH"); s != "" {
		list := filepath.SplitList(s)
		// TODO: multiple root dirs?
		defaultContext["rootDirectory"] = list[0]
	}
	return defaultContext
}

type Error struct {
	Type      string       `json:"type"`
	Message   string       `json:"message"`
	Backtrace []StackFrame `json:"backtrace"`
}

type Notice struct {
	Errors  []Error                `json:"errors"`
	Context map[string]interface{} `json:"context"`
	Env     map[string]interface{} `json:"environment"`
	Session map[string]interface{} `json:"session"`
	Params  map[string]interface{} `json:"params"`
}

func (n *Notice) String() string {
	if len(n.Errors) == 0 {
		return "Notice<no errors>"
	}
	e := n.Errors[0]
	return fmt.Sprintf("Notice<%s: %s>", e.Type, e.Message)
}

func NewNotice(e interface{}, req *http.Request, depth int) *Notice {
	notice := &Notice{
		Errors: []Error{{
			Type:      fmt.Sprintf("%T", e),
			Message:   fmt.Sprint(e),
			Backtrace: stack(depth),
		}},
		Context: map[string]interface{}{},
		Env:     map[string]interface{}{},
		Session: map[string]interface{}{},
		Params:  map[string]interface{}{},
	}

	for k, v := range getDefaultContext() {
		notice.Context[k] = v
	}

	if req != nil {
		notice.Context["url"] = req.URL.String()
		if ua := req.Header.Get("User-Agent"); ua != "" {
			notice.Context["userAgent"] = ua
		}

		for k, v := range req.Header {
			if len(v) == 1 {
				notice.Env[k] = v[0]
			} else {
				notice.Env[k] = v
			}
		}

		if err := req.ParseForm(); err == nil {
			for k, v := range req.Form {
				if len(v) == 1 {
					notice.Params[k] = v[0]
				} else {
					notice.Params[k] = v
				}
			}
		}
	}

	return notice
}
