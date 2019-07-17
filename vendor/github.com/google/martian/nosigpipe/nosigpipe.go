//+build !darwin !go1.9

package nosigpipe

import "net"

// IgnoreSIGPIPE prevents SIGPIPE from being raised on TCP sockets when remote hangs up
// See: https://github.com/golang/go/issues/17393. Do nothing for non Darwin
func IgnoreSIGPIPE(c net.Conn) {
}
