package base

import (
	"bytes"
	"encoding/json"
	"strings"
	"sync"
	"time"
	"useful-tools/module/logic/common"
	"useful-tools/pkg/wlog"

	"github.com/pkg/errors"
)

type Base struct {
	mt            sync.Mutex
	executing     bool
	proTemplate   string
	leftClickTime time.Time
	requestInfo   string
}

func (b *Base) ProTemplate() string {
	return b.proTemplate
}

func (b *Base) SetProTemplate(proTemplate string) {
	b.proTemplate = proTemplate
}

func (b *Base) IsExecuting() bool {
	return b.executing
}

func (b *Base) SetExecuting() {
	b.executing = true
}

func (b *Base) ResetExecuting() {
	b.executing = false
}

func (b *Base) IsExecutingError(err error) bool {
	return errors.Is(err, common.ExecutingError)
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
		wlog.Warm("json indent error: %v", err)
		return data
	}
	return strings.ReplaceAll(buf.String(), "\n", "\r\n")
}

func (b *Base) RequestInfo() string {
	return b.requestInfo
}

func (b *Base) SetRequestInfo(requestInfo string) {
	b.requestInfo = requestInfo
}
