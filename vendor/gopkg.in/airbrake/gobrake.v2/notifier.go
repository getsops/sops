package gobrake // import "gopkg.in/airbrake/gobrake.v2"

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const defaultAirbrakeHost = "https://airbrake.io"
const waitTimeout = 5 * time.Second
const httpStatusTooManyRequests = 429

var (
	errClosed      = errors.New("gobrake: notifier is closed")
	errRateLimited = errors.New("gobrake: rate limited")
)

var httpClient = &http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   15 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			ClientSessionCache: tls.NewLRUClientSessionCache(1024),
		},
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: 10 * time.Second,
	},
	Timeout: 10 * time.Second,
}

var buffers = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

type filter func(*Notice) *Notice

type Notifier struct {
	// http.Client that is used to interact with Airbrake API.
	Client *http.Client

	projectId       int64
	projectKey      string
	createNoticeURL string

	filters []filter

	wg       sync.WaitGroup
	noticeCh chan *Notice
	closed   chan struct{}
}

func NewNotifier(projectId int64, projectKey string) *Notifier {
	n := &Notifier{
		Client: httpClient,

		projectId:       projectId,
		projectKey:      projectKey,
		createNoticeURL: getCreateNoticeURL(defaultAirbrakeHost, projectId, projectKey),

		filters: []filter{noticeBacktraceFilter},

		noticeCh: make(chan *Notice, 1000),
		closed:   make(chan struct{}),
	}
	for i := 0; i < 10; i++ {
		go n.worker()
	}
	return n
}

// Sets Airbrake host name. Default is https://airbrake.io.
func (n *Notifier) SetHost(h string) {
	n.createNoticeURL = getCreateNoticeURL(h, n.projectId, n.projectKey)
}

// AddFilter adds filter that can modify or ignore notice.
func (n *Notifier) AddFilter(fn filter) {
	n.filters = append(n.filters, fn)
}

// Notify notifies Airbrake about the error.
func (n *Notifier) Notify(e interface{}, req *http.Request) {
	notice := n.Notice(e, req, 1)
	n.SendNoticeAsync(notice)
}

// Notice returns Aibrake notice created from error and request. depth
// determines which call frame to use when constructing backtrace.
func (n *Notifier) Notice(err interface{}, req *http.Request, depth int) *Notice {
	return NewNotice(err, req, depth+3)
}

type sendResponse struct {
	Id string `json:"id"`
}

// SendNotice sends notice to Airbrake.
func (n *Notifier) SendNotice(notice *Notice) (string, error) {
	for _, fn := range n.filters {
		notice = fn(notice)
		if notice == nil {
			// Notice is ignored.
			return "", nil
		}
	}

	buf := buffers.Get().(*bytes.Buffer)
	defer buffers.Put(buf)

	buf.Reset()
	if err := json.NewEncoder(buf).Encode(notice); err != nil {
		return "", err
	}

	resp, err := n.Client.Post(n.createNoticeURL, "application/json", buf)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	buf.Reset()
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusCreated {
		if resp.StatusCode == httpStatusTooManyRequests {
			return "", errRateLimited
		}
		err := fmt.Errorf("gobrake: got response status=%q, wanted 201 CREATED", resp.Status)
		return "", err
	}

	var sendResp sendResponse
	err = json.NewDecoder(buf).Decode(&sendResp)
	if err != nil {
		return "", err
	}

	return sendResp.Id, nil
}

func (n *Notifier) sendNotice(notice *Notice) {
	if _, err := n.SendNotice(notice); err != nil && err != errRateLimited {
		logger.Printf("gobrake failed reporting notice=%q: %s", notice, err)
	}
	n.wg.Done()
}

// SendNoticeAsync acts as SendNotice, but sends notice asynchronously
// and pending notices can be flushed with Flush.
func (n *Notifier) SendNoticeAsync(notice *Notice) {
	select {
	case <-n.closed:
		return
	default:
	}

	n.wg.Add(1)
	select {
	case n.noticeCh <- notice:
	default:
		n.wg.Done()
		logger.Printf(
			"notice=%q is ignored, because queue is full (len=%d)",
			notice, len(n.noticeCh),
		)
	}
}

func (n *Notifier) worker() {
	for {
		select {
		case notice := <-n.noticeCh:
			n.sendNotice(notice)
		case <-n.closed:
			select {
			case notice := <-n.noticeCh:
				n.sendNotice(notice)
			default:
				return
			}
		}
	}
}

// NotifyOnPanic notifies Airbrake about the panic and should be used
// with defer statement.
func (n *Notifier) NotifyOnPanic() {
	if v := recover(); v != nil {
		notice := n.Notice(v, nil, 3)
		n.SendNotice(notice)
		panic(v)
	}
}

// Flush waits for pending requests to finish.
func (n *Notifier) Flush() {
	n.waitTimeout(waitTimeout)
}

// Deprecated. Use CloseTimeout instead.
func (n *Notifier) WaitAndClose(timeout time.Duration) error {
	return n.CloseTimeout(timeout)
}

// CloseTimeout waits for pending requests to finish and then closes the notifier.
func (n *Notifier) CloseTimeout(timeout time.Duration) error {
	select {
	case <-n.closed:
	default:
		close(n.closed)
	}
	return n.waitTimeout(timeout)
}

func (n *Notifier) waitTimeout(timeout time.Duration) error {
	done := make(chan struct{})
	go func() {
		n.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("Wait timed out after %s", timeout)
	}
}

func (n *Notifier) Close() error {
	return n.CloseTimeout(waitTimeout)
}

func getCreateNoticeURL(host string, projectId int64, key string) string {
	return fmt.Sprintf(
		"%s/api/v3/projects/%d/notices?key=%s",
		host, projectId, key,
	)
}

func noticeBacktraceFilter(notice *Notice) *Notice {
	v, ok := notice.Context["rootDirectory"]
	if !ok {
		return notice
	}

	dir, ok := v.(string)
	if !ok {
		return notice
	}

	dir = filepath.Join(dir, "src")
	for i := range notice.Errors {
		replaceRootDirectory(notice.Errors[i].Backtrace, dir)
	}
	return notice
}

func replaceRootDirectory(backtrace []StackFrame, rootDir string) {
	for i := range backtrace {
		backtrace[i].File = strings.Replace(backtrace[i].File, rootDir, "[PROJECT_ROOT]", 1)
	}
}
