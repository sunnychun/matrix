package restful

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/ironzhang/matrix/codes"
	"github.com/ironzhang/matrix/tlog"
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
	log := tlog.Std().Sugar()

	e, ok := m.getEntry(r.URL.Path)
	if !ok {
		setErrorStatus(w, http.StatusNotFound, codes.Internal)
		log.Infow(http.StatusText(http.StatusNotFound), "method", r.Method, "path", r.URL.Path)
		return
	}
	h, ok := e.GetHandler(r.Method)
	if !ok {
		setErrorStatus(w, http.StatusMethodNotAllowed, codes.Internal)
		log.Infow(http.StatusText(http.StatusMethodNotAllowed), "method", r.Method, "path", r.URL.Path)
		return
	}
	m.serve(h, w, r)
}

func (m *ServeMux) serve(h *handler, w http.ResponseWriter, r *http.Request) {
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
	e := toJSONError(err)
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	fmt.Fprintln(w, e.Error())
}

func setErrorStatus(w http.ResponseWriter, status int, code codes.Code) {
	err := NewError(status, code)
	setError(w, err)
}
