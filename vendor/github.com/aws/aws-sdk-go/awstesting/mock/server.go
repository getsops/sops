package mock

import (
	"net/http"
	"net/http/httptest"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Session is a mock session which is used to hit the mock server
var Session = session.New(&aws.Config{
	DisableSSL: aws.Bool(true),
	Endpoint:   aws.String(server.URL[7:]),
})

// server is the mock server that simply writes a 200 status back to the client
var server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}))
