package version

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAIsNewerThanB(t *testing.T) {
	tests := []struct {
		A      string
		B      string
		Expect bool
	}{
		{"2.0.0", "1.0.0", true},
		{"1.1.0", "1.0.0", true},
		{"1.0.1", "1.0.0", true},
		{"v2.0.0", "v1.0.0", true},
		{"v1.1.0", "v1.0.0", true},
		{"v1.0.1", "v1.0.0", true},
		{"v1.0.0", "v1.0.0", false},
		{"v0.9.0", "v1.0.0", false},
		{"1.0.0", "1.0.0-alpha", true},
		{"0.9.0", "1.0.0", true}, // Special case, 1.x is always considered newer.
	}

	for _, test := range tests {
		result, err := AIsNewerThanB(test.A, test.B)
		if err != nil {
			t.Errorf("Error for A=%s, B=%s: %s", test.A, test.B, err)
			continue
		}
		if result != test.Expect {
			t.Errorf("Mismatch for A=%s, B=%s: got %v, expected %v", test.A, test.B, result, test.Expect)
		}
	}
}

func Test_releaseFetcher_LatestRelease(t *testing.T) {
	tests := []struct {
		name         string
		mockRedirect *mockServer
		mockAPI      *mockServer
		wantTag      string
		wantURL      string
		wantErr      bool
	}{
		{
			name: "RedirectSuccess",
			mockRedirect: &mockServer{
				redirectChain: []int{http.StatusMovedPermanently, http.StatusMovedPermanently},
				statusCode:    http.StatusFound,
				header:        http.Header{"Location": {"https://github.com/owner/repo/releases/tag/v2.0.0"}},
			},
			wantTag: "v2.0.0",
			wantURL: "https://github.com/owner/repo/releases/tag/v2.0.0",
		},
		{
			name: "APIFallbackSuccess",
			mockRedirect: &mockServer{
				statusCode: http.StatusNotFound,
			},
			mockAPI: &mockServer{
				statusCode: http.StatusOK,
				response:   `{"tag_name": "v1.0.0", "html_url": "https://github.com/owner/repo/releases/tag/v1.0.0"}`,
			},
			wantTag: "v1.0.0",
			wantURL: "https://github.com/owner/repo/releases/tag/v1.0.0",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := releaseFetcher{}

			var redirectServer *httptest.Server
			if test.mockRedirect != nil {
				redirectServer = test.mockRedirect.start()
				defer redirectServer.Close()
				f.endpoint = redirectServer.URL
			}

			var apiServer *httptest.Server
			if test.mockAPI != nil {
				apiServer = test.mockAPI.start()
				defer apiServer.Close()
				f.apiEndpoint = apiServer.URL
			}

			tag, url, err := f.LatestRelease("owner/repo")
			if test.wantErr {
				assert.Error(t, err, "Expected an error")
			} else {
				assert.NoError(t, err, "Unexpected error")
			}
			assert.Equal(t, test.wantTag, tag, "Incorrect release tag")
			assert.Equal(t, test.wantURL, url, "Incorrect release URL")
		})
	}
}

func Test_releaseFetcher_LatestReleaseUsingRedirect(t *testing.T) {
	tests := []struct {
		name       string
		mockServer mockServer
		wantTag    string
		wantURL    string
		wantErr    bool
	}{
		{
			name: "Success",
			mockServer: mockServer{
				// Include two redirects to ensure that the final redirect is used to determine the tag.
				redirectChain: []int{http.StatusMovedPermanently, http.StatusMovedPermanently},
				statusCode:    http.StatusFound,
				header:        http.Header{"Location": {"https://github.com/owner/repo/releases/tag/v1.0.0"}},
			},
			wantTag: "v1.0.0",
			wantURL: "https://github.com/owner/repo/releases/tag/v1.0.0",
		},
		{
			name: "RedirectMissingLocation",
			mockServer: mockServer{
				statusCode: http.StatusFound,
				header:     http.Header{},
			},
			wantErr: true,
		},
		{
			name: "UnexpectedStatusCode",
			mockServer: mockServer{
				statusCode: http.StatusOK,
			},
			wantErr: true,
		},
		{
			name: "TagNotFound",
			mockServer: mockServer{
				statusCode: http.StatusFound,
				header:     http.Header{"Location": {"https://github.com/owner/repo/releases"}},
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := test.mockServer.start()
			defer server.Close()

			f := releaseFetcher{
				endpoint: server.URL,
			}
			tag, url, err := f.LatestReleaseUsingRedirect("owner/repo")

			if test.wantErr {
				assert.Error(t, err, "Expected an error")
			} else {
				assert.NoError(t, err, "Unexpected error")
			}
			assert.Equal(t, test.wantTag, tag, "Incorrect release tag")
			assert.Equal(t, test.wantURL, url, "Incorrect release URL")
		})
	}
}

func Test_releaseFetcher_LatestReleaseUsingAPI(t *testing.T) {
	tests := []struct {
		name       string
		mockServer mockServer
		wantTag    string
		wantURL    string
		wantErr    bool
	}{
		{
			name: "Success",
			mockServer: mockServer{
				statusCode: http.StatusOK,
				response:   `{"tag_name": "v1.0.0", "html_url": "https://github.com/owner/repo/releases/tag/v1.0.0"}`,
			},
			wantTag: "v1.0.0",
			wantURL: "https://github.com/owner/repo/releases/tag/v1.0.0",
		},
		{
			name: "RequestError",
			mockServer: mockServer{
				statusCode: http.StatusInternalServerError,
				response:   "",
			},
			wantErr: true,
		},
		{
			name: "DecodeError",
			mockServer: mockServer{
				statusCode: http.StatusOK,
				response:   `{"invalid_json":}`,
			},
			wantErr: true,
		},
		{
			name: "NonOKStatusCode",
			mockServer: mockServer{
				statusCode: http.StatusNotFound,
				response:   "",
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := test.mockServer.start()
			defer server.Close()

			f := releaseFetcher{
				apiEndpoint: server.URL,
			}
			tag, url, err := f.LatestReleaseUsingAPI("owner/repo")

			if test.wantErr {
				assert.Error(t, err, "Expected an error")
			} else {
				assert.NoError(t, err, "Unexpected error")
			}
			assert.Equal(t, test.wantTag, tag, "Incorrect release tag")
			assert.Equal(t, test.wantURL, url, "Incorrect release URL")
		})
	}
}

type mockServer struct {
	statusCode    int
	header        http.Header
	response      string
	redirectChain []int
}

func (m *mockServer) start() *httptest.Server {
	redirectIndex := 0

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if redirectIndex < len(m.redirectChain) {
			redirectStatusCode := m.redirectChain[redirectIndex]
			redirectIndex++

			w.Header().Set("Location", "/redirected")
			w.WriteHeader(redirectStatusCode)
			return
		}

		for key, values := range m.header {
			w.Header()[key] = values
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(m.statusCode)
		fmt.Fprintln(w, m.response)
	}))
}
