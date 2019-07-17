//+build darwin,go1.9

package nosigpipe

import (
	"net"
	"syscall"

	"github.com/google/martian/v3/log"
)

// IgnoreSIGPIPE prevents SIGPIPE from being raised on TCP sockets when remote hangs up
// See: https://github.com/golang/go/issues/17393
func IgnoreSIGPIPE(c net.Conn) {
	if c == nil {
		return
	}
	s, ok := c.(syscall.Conn)
	if !ok {
		return
	}
	r, e := s.SyscallConn()
	if e != nil {
		log.Errorf("Failed to get SyscallConn: %s", e)
		return
	}
	e = r.Control(func(fd uintptr) {
		intfd := int(fd)
		if e := syscall.SetsockoptInt(intfd, syscall.SOL_SOCKET, syscall.SO_NOSIGPIPE, 1); e != nil {
			log.Errorf("Failed to set SO_NOSIGPIPE: %s", e)
		}
	})
	if e != nil {
		log.Errorf("Failed to set SO_NOSIGPIPE: %s", e)
	}
}
