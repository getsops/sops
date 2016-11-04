// +build js

package net

import (
	"errors"
	"syscall"

	"github.com/gopherjs/gopherjs/js"
)

func byteIndex(s string, c byte) int {
	return js.InternalObject(s).Call("indexOf", js.Global.Get("String").Call("fromCharCode", c)).Int()
}

func Listen(net, laddr string) (Listener, error) {
	panic(errors.New("network access is not supported by GopherJS"))
}

func (d *Dialer) Dial(network, address string) (Conn, error) {
	panic(errors.New("network access is not supported by GopherJS"))
}

func sysInit() {
}

func probeIPv4Stack() bool {
	return false
}

func probeIPv6Stack() (supportsIPv6, supportsIPv4map bool) {
	return false, false
}

func probeWindowsIPStack() (supportsVistaIP bool) {
	return false
}

func maxListenerBacklog() int {
	return syscall.SOMAXCONN
}
