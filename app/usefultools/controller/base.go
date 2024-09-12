package controller

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Base struct {
	mt            sync.Mutex
	leftClickTime time.Time
}

func (b *Base) DoubleClicked() (res bool) {
	b.mt.Lock()
	defer b.mt.Unlock()
	now := time.Now()
	if now.Sub(b.leftClickTime).Milliseconds() <= 800 {
		res = true
	} else {
		res = false
	}
	b.leftClickTime = now
	return
}

func (b *Base) FormatJson(data string) string {
	var buf bytes.Buffer
	err := json.Indent(&buf, []byte(data), "", "    ")
	if err != nil {
		logrus.Warnf("json indent error: %v", err)
		return data
	}
	return buf.String()
}
