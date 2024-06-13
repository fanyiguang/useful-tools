package useful

import (
	"time"
	"useful-tools/pkg/wlog"
)

func backGround() {
	go upgradeLoop()
}

func upgradeLoop() {
	ticker := time.NewTicker(10 * time.Minute)
	for {
		select {
		case <-ticker.C:
			err := upgrade()
			if err != nil {
				wlog.Warm("upgrade error: %v", err)
			}
		}
	}
}
