package walkUI

import (
	"sync"
	"testing"
	"time"
)

func TestNewApp(t *testing.T) {
	New()
}

func TestSyncOnce(t *testing.T) {
	var a sync.Once
	for i := 0; i < 10; i++ {
		go func() {
			a.Do(func() {
				time.Sleep(10 * time.Second)
			})
			t.Log("+++++++++", i)
		}()
	}
	select {}
}
