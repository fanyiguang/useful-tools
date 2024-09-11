package usefultools

import (
	"github.com/sirupsen/logrus"
	"time"
)

func backGround() {
	go upgradeLoop()
}

func upgradeLoop() {
	ticker := time.NewTicker(2 * time.Minute)
	for {
		select {
		case <-ticker.C:
			err := upgrade()
			if err != nil {
				logrus.Warnf("upgrade error: %v", err)
			}
		}
	}
}
