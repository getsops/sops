// Copyright 2017 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mobile

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/api"
	"github.com/google/martian/v3/cors"
	"github.com/google/martian/v3/cybervillains"
	"github.com/google/martian/v3/fifo"
	"github.com/google/martian/v3/har"
	"github.com/google/martian/v3/httpspec"
	mlog "github.com/google/martian/v3/log"
	"github.com/google/martian/v3/marbl"
	"github.com/google/martian/v3/martianhttp"
	"github.com/google/martian/v3/mitm"
	"github.com/google/martian/v3/servemux"
	"github.com/google/martian/v3/trafficshape"
	"github.com/google/martian/v3/verify"

	// side-effect importing to register with JSON API
	_ "github.com/google/martian/v3/body"
	_ "github.com/google/martian/v3/cookie"
	_ "github.com/google/martian/v3/failure"
	_ "github.com/google/martian/v3/header"
	_ "github.com/google/martian/v3/martianurl"
	_ "github.com/google/martian/v3/method"
	_ "github.com/google/martian/v3/pingback"
	_ "github.com/google/martian/v3/port"
	_ "github.com/google/martian/v3/priority"
	_ "github.com/google/martian/v3/querystring"
	_ "github.com/google/martian/v3/skip"
	_ "github.com/google/martian/v3/stash"
	_ "github.com/google/martian/v3/static"
	_ "github.com/google/martian/v3/status"
)

// Martian is a wrapper for the initialized Martian proxy
type Martian struct {
	proxy          *martian.Proxy
	listener       net.Listener
	apiListener    net.Listener
	mux            *http.ServeMux
	started        bool
	HARLogging     bool
	TrafficPort    int
	TrafficShaping bool
	APIPort        int
	APIOverTLS     bool
	BindLocalhost  bool
	Cert           string
	Key            string
	AllowCORS      bool
	RoundTripper   *http.Transport
}

// EnableCybervillains configures Martian to use the Cybervillians certificate.
func (m *Martian) EnableCybervillains() {
	m.Cert = cybervillains.Cert
	m.Key = cybervillains.Key
}

// NewProxy creates a new Martian struct for configuring and starting a martian.
func NewProxy() *Martian {
	return &Martian{}
}

// Start starts the proxy given the configured values of the Martian struct.
func (m *Martian) Start() {
	var err error
	m.listener, err = net.Listen("tcp", m.bindAddress(m.TrafficPort))
	if err != nil {
		log.Fatal(err)
	}

	mlog.Debugf("mobile: started listener on: %v", m.listener.Addr())
	m.proxy = martian.NewProxy()
	m.mux = http.NewServeMux()

	if m.Cert != "" && m.Key != "" {
		tlsc, err := tls.X509KeyPair([]byte(m.Cert), []byte(m.Key))
		if err != nil {
			log.Fatal(err)
		}

		mlog.Debugf("mobile: loaded cert and key")

		x509c, err := x509.ParseCertificate(tlsc.Certificate[0])
		if err != nil {
			log.Fatal(err)
		}

		mlog.Debugf("mobile: parsed cert")

		mc, err := mitm.NewConfig(x509c, tlsc.PrivateKey)
		if err != nil {
			log.Fatal(err)
		}

		mc.SetValidity(12 * time.Hour)
		mc.SetOrganization("Martian Proxy")

		m.proxy.SetMITM(mc)

		if m.RoundTripper != nil {
			m.proxy.SetRoundTripper(m.RoundTripper)
		}
		m.handle("/authority.cer", martianhttp.NewAuthorityHandler(x509c))
	}

	// Enable Traffic shaping if requested
	if m.TrafficShaping {
		tsl := trafficshape.NewListener(m.listener)
		tsh := trafficshape.NewHandler(tsl)
		m.handle("/shape-traffic", tsh)
		m.listener = tsl
	}

	// Forward traffic that pattern matches in m.mux before applying
	// httpspec modifiers (via modifier, specifically)
	topg := fifo.NewGroup()
	apif := servemux.NewFilter(m.mux)
	apif.SetRequestModifier(api.NewForwarder("", m.APIPort))
	topg.AddRequestModifier(apif)

	stack, fg := httpspec.NewStack("martian.mobile")
	topg.AddRequestModifier(stack)
	topg.AddResponseModifier(stack)

	m.proxy.SetRequestModifier(topg)
	m.proxy.SetResponseModifier(topg)

	if m.HARLogging {
		// add HAR logger for unmodified logs.
		uhl := har.NewLogger()
		uhmuxf := servemux.NewFilter(m.mux)
		uhmuxf.RequestWhenFalse(uhl)
		uhmuxf.ResponseWhenFalse(uhl)
		fg.AddRequestModifier(uhmuxf)
		fg.AddResponseModifier(uhmuxf)

		// add HAR logger
		hl := har.NewLogger()
		hmuxf := servemux.NewFilter(m.mux)
		hmuxf.RequestWhenFalse(hl)
		hmuxf.ResponseWhenFalse(hl)
		stack.AddRequestModifier(hmuxf)
		stack.AddResponseModifier(hmuxf)

		// Retrieve Unmodified HAR logs
		m.handle("/logs/original", har.NewExportHandler(uhl))
		m.handle("/logs/original/reset", har.NewResetHandler(uhl))

		// Retrieve HAR logs
		m.handle("/logs", har.NewExportHandler(hl))
		m.handle("/logs/reset", har.NewResetHandler(hl))
	}

	lsh := marbl.NewHandler()
	// retrieve binary marbl logs
	m.handle("/binlogs", lsh)

	lsm := marbl.NewModifier(lsh)
	muxf := servemux.NewFilter(m.mux)
	muxf.RequestWhenFalse(lsm)
	muxf.ResponseWhenFalse(lsm)
	stack.AddRequestModifier(muxf)
	stack.AddResponseModifier(muxf)

	mod := martianhttp.NewModifier()
	fg.AddRequestModifier(mod)
	fg.AddResponseModifier(mod)

	// Proxy specific handlers.
	// These handlers take precendence over proxy traffic and will not be intercepted.

	// Update modifiers.
	m.handle("/configure", mod)

	// Verify assertions.
	vh := verify.NewHandler()
	vh.SetRequestVerifier(mod)
	vh.SetResponseVerifier(mod)

	m.handle("/verify", vh)

	// Reset verifications.
	rh := verify.NewResetHandler()
	rh.SetRequestVerifier(mod)
	rh.SetResponseVerifier(mod)
	m.handle("/verify/reset", rh)

	mlog.Infof("mobile: starting Martian proxy on listener")
	go m.proxy.Serve(m.listener)

	// start the API server
	apiAddr := m.bindAddress(m.APIPort)
	m.apiListener, err = net.Listen("tcp", apiAddr)
	if err != nil {
		log.Fatal(err)
	}
	if m.APIOverTLS {
		if m.Cert == "" || m.Key == "" {
			log.Fatal("mobile: APIOverTLS cannot be true without valid cert and key")
		}

		cerfile, err := ioutil.TempFile("", "martian-api.cert")
		if err != nil {
			log.Fatal(err)
		}

		keyfile, err := ioutil.TempFile("", "martian-api.key")
		if err != nil {
			log.Fatal(err)
		}

		if _, err := cerfile.Write([]byte(m.Cert)); err != nil {
			log.Fatal(err)
		}

		if _, err := keyfile.Write([]byte(m.Key)); err != nil {
			log.Fatal(err)
		}

		go func() {
			http.ServeTLS(m.apiListener, m.mux, cerfile.Name(), keyfile.Name())
			defer os.Remove(cerfile.Name())
			defer os.Remove(keyfile.Name())
		}()

		mlog.Infof("mobile: proxy API started on %s over TLS", apiAddr)
	} else {
		go http.Serve(m.apiListener, m.mux)
		mlog.Infof("mobile: proxy API started on %s", apiAddr)
	}

	m.started = true
}

// IsStarted returns true if the proxy has finished starting.
func (m *Martian) IsStarted() bool {
	return m.started
}

// Shutdown tells the Proxy to close. This function returns immediately, though
// there may still be connection threads hanging around until they time out
// depending on how the OS manages them.
func (m *Martian) Shutdown() {
	mlog.Infof("mobile: shutting down proxy")
	m.listener.Close()
	m.apiListener.Close()
	m.proxy.Close()
	m.started = false
	mlog.Infof("mobile: proxy shut down")
}

// SetLogLevel sets the Martian log level (Silent = 0, Error, Info, Debug), controlling which Martian
// log calls are displayed in the console
func SetLogLevel(l int) {
	mlog.SetLevel(l)
}

func (m *Martian) handle(pattern string, handler http.Handler) {
	if m.AllowCORS {
		handler = cors.NewHandler(handler)
	}
	m.mux.Handle(pattern, handler)
	mlog.Infof("mobile: handler registered for %s", pattern)

	lhp := path.Join(fmt.Sprintf("localhost:%d", m.APIPort), pattern)
	m.mux.Handle(lhp, handler)
	mlog.Infof("mobile: handler registered for %s", lhp)
}

func (m *Martian) bindAddress(port int) string {
	if m.BindLocalhost {
		return fmt.Sprintf("[::1]:%d", port)
	}
	return fmt.Sprintf(":%d", port)
}
