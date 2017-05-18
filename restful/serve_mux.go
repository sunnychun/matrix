package restful

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/ironzhang/matrix/codes"
	"github.com/ironzhang/matrix/context-value"
	"github.com/ironzhang/matrix/tlog"
	"github.com/ironzhang/matrix/uuid"
)

const contentType = "application/json"

type entry struct {
	mu sync.RWMutex
	m  map[string]*handler
}

func (e *entry) AddHandler(meth string, h *handler) {
	meth = strings.ToUpper(meth)
	e.mu.Lock()
	e.m[meth] = h
	e.mu.Unlock()
}

func (e *entry) GetHandler(meth string) (*handler, bool) {
	meth = strings.ToUpper(meth)
	e.mu.RLock()
	h, ok := e.m[meth]
	e.mu.RUnlock()
	return h, ok
}

func NewServeMux() *ServeMux {
	return &ServeMux{
		verbose: 1,
		entrys:  make(map[string]*entry),
	}
}

type ServeMux struct {
	verbose int
	mu      sync.RWMutex
	entrys  map[string]*entry
}

func (m *ServeMux) Delete(pat string, h interface{}) error {
	return m.Add("DELETE", pat, h)
}

func (m *ServeMux) Get(pat string, h interface{}) error {
	return m.Add("GET", pat, h)
}

func (m *ServeMux) Head(pat string, h interface{}) error {
	return m.Add("HEAD", pat, h)
}

func (m *ServeMux) Options(pat string, h interface{}) error {
	return m.Add("OPTIONS", pat, h)
}

func (m *ServeMux) Patch(pat string, h interface{}) error {
	return m.Add("PATCH", pat, h)
}

func (m *ServeMux) Post(pat string, h interface{}) error {
	return m.Add("POST", pat, h)
}

func (m *ServeMux) Put(pat string, h interface{}) error {
	return m.Add("PUT", pat, h)
}

func (m *ServeMux) Add(meth, pat string, i interface{}) error {
	h, err := parseHandler(i)
	if err != nil {
		return fmt.Errorf("parse handler: %v", err)
	}
	m.addEntry(pat).AddHandler(meth, h)
	return nil
}

func (m *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context_value.WithTraceId(context.Background(), getTraceId(r.Header))
	if err := m.serveHTTP(ctx, w, r); err != nil {
		setError(w, err)
	}
}

func (m *ServeMux) serveHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	log := tlog.WithContext(ctx).Sugar()
	e, ok := m.getEntry(r.URL.Path)
	if !ok {
		log.Infow(http.StatusText(http.StatusNotFound), "method", r.Method, "path", r.URL.Path)
		return Errorf(http.StatusNotFound, codes.NotFound, "page(%s) not found", r.URL.Path)
	}
	h, ok := e.GetHandler(r.Method)
	if !ok {
		log.Infow(http.StatusText(http.StatusMethodNotAllowed), "method", r.Method, "path", r.URL.Path)
		return Errorf(http.StatusMethodNotAllowed, codes.NotAllowed, "method(%s) not allowed", r.Method)
	}
	return m.serve(ctx, h, w, r)
}

func (m *ServeMux) serve(ctx context.Context, h *handler, w http.ResponseWriter, r *http.Request) error {
	log := tlog.WithContext(ctx).Sugar()

	// check Content-Type
	if v := r.Header.Get("Content-Type"); v != "" && v != contentType {
		cause := fmt.Sprintf("Content-Type(%s) not %s", v, contentType)
		log.Infow(cause, "method", r.Method, "path", r.URL.Path)
		return Errorf(http.StatusBadRequest, codes.InvalidHeader, cause)
	}

	var err error
	in1 := newReflectValue(h.in1Type)
	in2 := newReflectValue(h.in2Type)

	// Decode
	if !isNilInterface(h.in1Type) {
		if err = json.NewDecoder(r.Body).Decode(in1.Interface()); err != nil {
			log.Infow("decode fail", "error", err, "method", r.Method, "path", r.URL.Path)
			return Errorf(http.StatusBadRequest, codes.DecodeFail, err.Error())
		}
	}

	// Handle
	if err = h.Handle(ctx, in1, in2); err != nil {
		log.Infow("handle fail", "error", err, "method", r.Method, "path", r.URL.Path)
		return err
	}

	// Encode
	if !isNilInterface(h.in2Type) {
		if err = json.NewEncoder(w).Encode(in2.Interface()); err != nil {
			log.Errorw("encode fail", "error", err, "method", r.Method, "path", r.URL.Path)
			return Errorf(http.StatusInternalServerError, codes.EncodeFail, err.Error())
		}
	}

	return nil
}

func (m *ServeMux) addEntry(pat string) *entry {
	m.mu.Lock()
	e, ok := m.entrys[pat]
	if !ok {
		e = &entry{m: make(map[string]*handler)}
		m.entrys[pat] = e
	}
	m.mu.Unlock()
	return e
}

func (m *ServeMux) getEntry(pat string) (*entry, bool) {
	m.mu.RLock()
	e, ok := m.entrys[pat]
	m.mu.RUnlock()
	return e, ok
}

func setError(w http.ResponseWriter, err error) {
	status := http.StatusBadRequest
	if te, ok := err.(HTTPStatus); ok {
		status = te.HTTPStatus()
	}
	e := toRPCError(err)
	w.WriteHeader(status)
	fmt.Fprintln(w, e.Error())
}

func getTraceId(h http.Header) string {
	if v := h.Get("X-Trace-Id"); v != "" {
		return v
	}
	return uuid.New().String()
}
