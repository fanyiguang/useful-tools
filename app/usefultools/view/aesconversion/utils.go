package aesconversion

import (
	"fmt"
	"fyne.io/fyne/v2"
	"strings"
	"time"
	"useful-tools/app/usefultools/view/constant"
)

func getKeyGroupList() ([]string, map[string]string, map[string]string) {
	var names []string
	key := make(map[string]string)
	iv := make(map[string]string)
	keyList := fyne.CurrentApp().Preferences().StringList("aes-key-list")
	fmt.Println(keyList)
	for _, k := range keyList {
		keyInfo := strings.Split(k, ":")
		if len(keyInfo) != 3 {
			continue
		}
		names = append(names, keyInfo[0])
		key[keyInfo[0]] = keyInfo[1]
		iv[keyInfo[0]] = keyInfo[2]
	}
	return names, key, iv
}

func setKeyGroupList(name, key, iv string) (changed bool) {
	if fyne.CurrentApp().Preferences().Bool(constant.NavStatePreferenceSaveAesKey) {
		if name == "" {
			name = fmt.Sprintf("%s-%v", time.Now().Format(`2006-01-02`), time.Now().UnixMilli())
		}
		found := false
		content := fmt.Sprintf("%s:%s:%s", name, key, iv)
		keyList := fyne.CurrentApp().Preferences().StringList("aes-key-list")
		for i, k := range keyList {
			keyInfo := strings.Split(k, ":")
			if len(keyInfo) != 3 {
				continue
			}
			if keyInfo[0] == name {
				found = true
				keyList = append(keyList[:i], keyList[i+1:]...)
				keyList = append([]string{content}, keyList...)
				changed = true
				break
			}
		}

		if !found {
			changed = true
			keyList = append(keyList, content)
		}

		if len(keyList) > 22 {
			keyList = keyList[0:22]
		}

		fyne.CurrentApp().Preferences().SetStringList("aes-key-list", keyList)
	}
	return
}
