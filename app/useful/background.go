package useful

import (
	"os"
	"time"
	"useful-tools/pkg/wlog"
)

func backGround() {
	go upgradeLoop()
}

func upgradeLoop() {
	os.Rename()
	ticker := time.NewTicker(1 * time.Minute)
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
