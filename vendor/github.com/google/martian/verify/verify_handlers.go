// Copyright 2015 Google Inc. All rights reserved.
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

package verify

import (
	"encoding/json"
	"net/http"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/log"
)

// Handler is an http.Handler that returns the request and response
// verifications of reqv and resv as JSON.
type Handler struct {
	reqv RequestVerifier
	resv ResponseVerifier
}

// ResetHandler is an http.Handler that resets the request and response
// verifications of reqv and resv.
type ResetHandler struct {
	reqv RequestVerifier
	resv ResponseVerifier
}

type verifyResponse struct {
	Errors []verifyError `json:"errors"`
}

type verifyError struct {
	Message string `json:"message"`
}

// NewHandler returns an http.Handler for requesting the verification
// error status.
func NewHandler() *Handler {
	return &Handler{}
}

// NewResetHandler returns an http.Handler for reseting the verification error
// status.
func NewResetHandler() *ResetHandler {
	return &ResetHandler{}
}

// SetRequestVerifier sets the RequestVerifier to verify.
func (h *Handler) SetRequestVerifier(reqv RequestVerifier) {
	h.reqv = reqv
}

// SetResponseVerifier sets the ResponseVerifier to verify.
func (h *Handler) SetResponseVerifier(resv ResponseVerifier) {
	h.resv = resv
}

// ServeHTTP writes out a JSON response containing a list of verification
// errors that occurred during the requests and responses sent to the proxy.
func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	if req.Method != "GET" {
		rw.Header().Set("Allow", "GET")
		rw.WriteHeader(405)
		log.Errorf("verify: invalid request method: %s", req.Method)
		return
	}

	vres := &verifyResponse{
		Errors: make([]verifyError, 0),
	}

	if h.reqv != nil {
		if err := h.reqv.VerifyRequests(); err != nil {
			appendError(vres, err)
		}
	}
	if h.resv != nil {
		if err := h.resv.VerifyResponses(); err != nil {
			appendError(vres, err)
		}
	}

	json.NewEncoder(rw).Encode(vres)
}

func appendError(vres *verifyResponse, err error) {
	merr, ok := err.(*martian.MultiError)
	if !ok {
		vres.Errors = append(vres.Errors, verifyError{Message: err.Error()})
		return
	}

	for _, err := range merr.Errors() {
		vres.Errors = append(vres.Errors, verifyError{Message: err.Error()})
	}
}

// SetRequestVerifier sets the RequestVerifier to reset.
func (h *ResetHandler) SetRequestVerifier(reqv RequestVerifier) {
	h.reqv = reqv
}

// SetResponseVerifier sets the ResponseVerifier to reset.
func (h *ResetHandler) SetResponseVerifier(resv ResponseVerifier) {
	h.resv = resv
}

// ServeHTTP resets the verifier for the given ID so that it may
// be run again.
func (h *ResetHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		rw.Header().Set("Allow", "POST")
		rw.WriteHeader(405)
		log.Errorf("verify: invalid request method: %s", req.Method)
		return
	}

	if h.reqv != nil {
		h.reqv.ResetRequestVerifications()
	}
	if h.resv != nil {
		h.resv.ResetResponseVerifications()
	}

	rw.WriteHeader(204)
}
