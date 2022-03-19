package base

import (
	"strings"
	"sync"
)

type Serializable struct {
	sync.Mutex
	concurrentParserParams string
}

func NewSerializable() *Serializable {
	return &Serializable{}
}

func (s *Serializable) Set(args ...string) string {
	s.Lock()
	defer s.Unlock()
	joinStr := strings.Join(args, ",")
	s.concurrentParserParams = joinStr
	return joinStr
}

func (s *Serializable) Get() string {
	s.Lock()
	defer s.Unlock()
	return s.concurrentParserParams
}

func (s *Serializable) Equal(str string) bool {
	return s.concurrentParserParams == str
}
