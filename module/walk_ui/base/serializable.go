package base

import (
	"strings"
	"sync"
)

type Persistence struct {
	sync.Mutex
	latestParams string
}

func NewPersistence() *Persistence {
	return &Persistence{}
}

func (s *Persistence) SetLatestParams(args ...string) string {
	s.Lock()
	defer s.Unlock()
	s.latestParams = strings.Join(args, ",")
	return s.latestParams
}

func (s *Persistence) GetLatestParams() string {
	s.Lock()
	defer s.Unlock()
	return s.latestParams
}

func (s *Persistence) Equal(str string) bool {
	return s.latestParams == str
}
