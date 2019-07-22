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

package static

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/proxyutil"
)

func Test404WhenExplictlyMappedFileDoesNotExist(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "test_static_modifier_explicit_path_mapping_")
	if err != nil {
		t.Fatalf("ioutil.TempDir(): got %v, want no error", err)
	}

	//if err := os.MkdirAll(path.Join(tmpdir, "explicit/path"), 0777); err != nil {
	//	t.Fatalf("os.Mkdir(): got %v, want no error", err)
	//}

	//if err := ioutil.WriteFile(path.Join(tmpdir, "explicit/path", "sfmtest.txt"), []byte("test file"), 0777); err != nil {
	//	t.Fatalf("ioutil.WriteFile(): got %v, want no error", err)
	//}

	req, err := http.NewRequest("GET", "/sfmtest.txt", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("TestContext(): got %v, want no error", err)
	}
	defer remove()

	res := proxyutil.NewResponse(http.StatusOK, nil, req)

	mod := NewModifier(tmpdir)
	if err := mod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	mod.SetExplicitPathMappings(map[string]string{"/sfmtest.txt": "/explicit/path/sfmtest.txt"})

	if err := mod.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	if got, want := res.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("res.StatusCode: got %v, want %v", got, want)
	}
}

func TestFileExistsInBothExplictlyMappedPathAndInferredPath(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "test_static_modifier_explicit_path_mapping_")
	if err != nil {
		t.Fatalf("ioutil.TempDir(): got %v, want no error", err)
	}

	if err := os.MkdirAll(path.Join(tmpdir, "explicit/path"), 0777); err != nil {
		t.Fatalf("os.Mkdir(): got %v, want no error", err)
	}

	if err := ioutil.WriteFile(path.Join(tmpdir, "sfmtest.txt"), []byte("dont return"), 0777); err != nil {
		t.Fatalf("ioutil.WriteFile(): got %v, want no error", err)
	}

	if err := ioutil.WriteFile(path.Join(tmpdir, "explicit/path", "sfmtest.txt"), []byte("target"), 0777); err != nil {
		t.Fatalf("ioutil.WriteFile(): got %v, want no error", err)
	}

	req, err := http.NewRequest("GET", "/sfmtest.txt", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("TestContext(): got %v, want no error", err)
	}
	defer remove()

	res := proxyutil.NewResponse(http.StatusOK, nil, req)

	mod := NewModifier(tmpdir)
	if err := mod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	mod.SetExplicitPathMappings(map[string]string{"/sfmtest.txt": "/explicit/path/sfmtest.txt"})

	if err := mod.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	if got, want := res.Header.Get("Content-Type"), "text/plain; charset=utf-8"; got != want {
		t.Errorf("res.Header.Get('Content-Type'): got %v, want %v", got, want)
	}

	got, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(): got %v, want no error", err)
	}
	res.Body.Close()

	if want := []byte("target"); !bytes.Equal(got, want) {
		t.Errorf("res.Body: got %q, want %q", got, want)
	}

	if got, want := res.ContentLength, int64(len("target")); got != want {
		t.Errorf("res.ContentLength: got %v, want %v", got, want)
	}
}

func TestStaticModifierExplicitPathMapping(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "test_static_modifier_explicit_path_mapping_")
	if err != nil {
		t.Fatalf("ioutil.TempDir(): got %v, want no error", err)
	}

	if err := os.MkdirAll(path.Join(tmpdir, "explicit/path"), 0777); err != nil {
		t.Fatalf("os.Mkdir(): got %v, want no error", err)
	}

	if err := ioutil.WriteFile(path.Join(tmpdir, "explicit/path", "sfmtest.txt"), []byte("test file"), 0777); err != nil {
		t.Fatalf("ioutil.WriteFile(): got %v, want no error", err)
	}

	req, err := http.NewRequest("GET", "/sfmtest.txt", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("TestContext(): got %v, want no error", err)
	}
	defer remove()

	res := proxyutil.NewResponse(http.StatusOK, nil, req)

	mod := NewModifier(tmpdir)
	if err := mod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	mod.SetExplicitPathMappings(map[string]string{"/sfmtest.txt": "/explicit/path/sfmtest.txt"})

	if err := mod.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	if got, want := res.Header.Get("Content-Type"), "text/plain; charset=utf-8"; got != want {
		t.Errorf("res.Header.Get('Content-Type'): got %v, want %v", got, want)
	}

	got, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(): got %v, want no error", err)
	}
	res.Body.Close()

	if want := []byte("test file"); !bytes.Equal(got, want) {
		t.Errorf("res.Body: got %q, want %q", got, want)
	}

	if got, want := res.ContentLength, int64(len("test file")); got != want {
		t.Errorf("res.ContentLength: got %v, want %v", got, want)
	}
}

func TestStaticModifierOnRequest(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "test_static_modifier_on_request_")
	if err != nil {
		t.Fatalf("ioutil.TempDir(): got %v, want no error", err)
	}

	if err := ioutil.WriteFile(path.Join(tmpdir, "sfmtest.txt"), []byte("test file"), 0777); err != nil {
		t.Fatalf("ioutil.WriteFile(): got %v, want no error", err)
	}

	req, err := http.NewRequest("GET", "/sfmtest.txt", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("TestContext(): got %v, want no error", err)
	}
	defer remove()

	res := proxyutil.NewResponse(http.StatusOK, nil, req)

	mod := NewModifier(tmpdir)
	if err := mod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)

	}
	if err := mod.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	if got, want := res.Header.Get("Content-Type"), "text/plain; charset=utf-8"; got != want {
		t.Errorf("res.Header.Get('Content-Type'): got %v, want %v", got, want)
	}

	got, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(): got %v, want no error", err)
	}
	res.Body.Close()

	if want := []byte("test file"); !bytes.Equal(got, want) {
		t.Errorf("res.Body: got %q, want %q", got, want)
	}

	if got, want := res.ContentLength, int64(len("test file")); got != want {
		t.Errorf("res.ContentLength: got %v, want %v", got, want)
	}
}

func TestRequestOverHTTPS(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "test_static_modifier_on_request_")
	if err != nil {
		t.Fatalf("ioutil.TempDir(): got %v, want no error", err)
	}

	if err := ioutil.WriteFile(path.Join(tmpdir, "sfmtest.txt"), []byte("test file"), 0777); err != nil {
		t.Fatalf("ioutil.WriteFile(): got %v, want no error", err)
	}

	req, err := http.NewRequest("GET", "/sfmtest.txt", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	req.URL.Scheme = "https"

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("TestContext(): got %v, want no error", err)
	}
	defer remove()

	res := proxyutil.NewResponse(http.StatusOK, nil, req)

	mod := NewModifier(tmpdir)
	if err := mod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)

	}
	if err := mod.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	if got, want := res.Header.Get("Content-Type"), "text/plain; charset=utf-8"; got != want {
		t.Errorf("res.Header.Get('Content-Type'): got %v, want %v", got, want)
	}

	got, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(): got %v, want no error", err)
	}
	res.Body.Close()

	if want := []byte("test file"); !bytes.Equal(got, want) {
		t.Errorf("res.Body: got %q, want %q", got, want)
	}

	if got, want := res.ContentLength, int64(len("test file")); got != want {
		t.Errorf("res.ContentLength: got %v, want %v", got, want)
	}

}

func TestModifierFromJSON(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "test_static_modifier_on_request_")
	if err != nil {
		t.Fatalf("ioutil.TempDir(): got %v, want no error", err)
	}

	tmpdir2 := path.Join(tmpdir, "subdir")
	err = os.Mkdir(tmpdir2, 0777)
	if err != nil {
		t.Fatalf("ioutil.TempDir(): got %v, want no error", err)
	}

	if err := ioutil.WriteFile(path.Join(tmpdir, "sfmtest.txt"), []byte("test file"), 0777); err != nil {
		t.Fatalf("ioutil.WriteFile(): got %v, want no error", err)
	}

	if err := ioutil.WriteFile(path.Join(tmpdir2, "sfmtest.txt"), []byte("test file2"), 0777); err != nil {
		t.Fatalf("ioutil.WriteFile(): got %v, want no error", err)
	}

	msg := []byte(fmt.Sprintf(`{
		"static.Modifier": {
			"scope": ["request", "response"],
			"explicitPaths": {"/foo/bar.baz": "/subdir/sfmtest.txt"},
			"rootPath": %q
		}
	}`, tmpdir))

	r, err := parse.FromJSON(msg)
	if err != nil {
		t.Fatalf("parse.FromJSON(): got %v, want no error", err)
	}

	reqmod := r.RequestModifier()
	if reqmod == nil {
		t.Fatal("reqmod: got nil, want not nil")
	}

	resmod := r.ResponseModifier()
	if resmod == nil {
		t.Fatal("resmod: got nil, want not nil")
	}

	req, err := http.NewRequest("GET", "/sfmtest.txt", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("TestContext(): got %v, want no error", err)
	}
	defer remove()

	res := proxyutil.NewResponse(http.StatusOK, nil, req)

	if err := reqmod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if err := resmod.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	if got, want := res.Header.Get("Content-Type"), "text/plain; charset=utf-8"; got != want {
		t.Errorf("res.Header.Get('Content-Type'): got %v, want %v", got, want)
	}

	got, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(): got %v, want no error", err)
	}
	res.Body.Close()

	if want := []byte("test file"); !bytes.Equal(got, want) {
		t.Errorf("res.Body: got %q, want %q", got, want)
	}

	if got, want := res.ContentLength, int64(len("test file")); got != want {
		t.Errorf("res.ContentLength: got %v, want %v", got, want)
	}

	req, err = http.NewRequest("GET", "/foo/bar.baz", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	_, remove, err = martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("TestContext(): got %v, want no error", err)
	}
	defer remove()

	res = proxyutil.NewResponse(http.StatusOK, nil, req)

	if err := reqmod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if err := resmod.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	if got, want := res.Header.Get("Content-Type"), "text/plain; charset=utf-8"; got != want {
		t.Errorf("res.Header.Get('Content-Type'): got %v, want %v", got, want)
	}

	got, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(): got %v, want no error", err)
	}
	res.Body.Close()

	if want := []byte("test file2"); !bytes.Equal(got, want) {
		t.Errorf("res.Body: got %q, want %q", got, want)
	}

	if got, want := res.ContentLength, int64(len("test file2")); got != want {
		t.Errorf("res.ContentLength: got %v, want %v", got, want)
	}

}

func TestStaticModifierSingleRangeRequest(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "test_static_modifier_on_request_")
	if err != nil {
		t.Fatalf("ioutil.TempDir(): got %v, want no error", err)
	}
	mod := NewModifier(tmpdir)

	if err := ioutil.WriteFile(path.Join(tmpdir, "sfmtest.txt"), []byte("0123456789"), 0777); err != nil {
		t.Fatalf("ioutil.WriteFile(): got %v, want no error", err)
	}

	req, err := http.NewRequest("GET", "/sfmtest.txt", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}
	req.Header.Set("Range", "bytes=1-4")

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("TestContext(): got %v, want no error", err)
	}
	defer remove()

	if err := mod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)

	}

	res := proxyutil.NewResponse(http.StatusOK, nil, req)
	if err := mod.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	if got, want := res.StatusCode, http.StatusPartialContent; got != want {
		t.Errorf("res.Status: got %v, want %v", got, want)
	}

	if got, want := res.ContentLength, int64(len([]byte("1234"))); got != want {
		t.Errorf("res.ContentLength: got %d, want %d", got, want)
	}

	if got, want := res.Header.Get("Content-Range"), "bytes 1-4/10"; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "Content-Encoding", got, want)
	}

	got, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(): got %v, want no error", err)
	}
	res.Body.Close()

	if want := []byte("1234"); !bytes.Equal(got, want) {
		t.Errorf("res.Body: got %q, want %q", got, want)
	}

	if got, want := res.Header.Get("Content-Type"), "text/plain; charset=utf-8"; got != want {
		t.Errorf("res.Header.Get('Content-Type'): got %v, want %v", got, want)
	}
}
