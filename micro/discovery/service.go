package discovery

import (
	"sort"
	"sync"
)

type Service interface {
	Name() string
	Addrs() []string
}

func newService(name string, refreshs ...Refresh) *service {
	return &service{
		name:     name,
		refreshs: refreshs,
		done:     make(chan struct{}),
	}
}

type service struct {
	name     string
	refreshs []Refresh
	done     chan struct{}
	ok       <-chan struct{}
	mu       sync.RWMutex
	addrs    []string
}

func (s *service) Name() string {
	return s.name
}

func (s *service) Addrs() []string {
	s.mu.RLock()
	addrs := s.addrs
	s.mu.RUnlock()
	return addrs
}

func (s *service) SetAddrs(addrs []string) {
	s.mu.Lock()
	s.addrs = addrs
	s.mu.Unlock()
}

func (s *service) Refresh(kvs map[string][]byte) {
	addrs := make([]string, 0, len(kvs))
	for k, _ := range kvs {
		addrs = append(addrs, k)
	}
	sort.Strings(addrs)
	s.SetAddrs(addrs)

	for _, refresh := range s.refreshs {
		refresh(addrs)
	}
}

func (s *service) Unwatch() {
	close(s.done)
	<-s.ok
}
