package preferencekey

import (
	"fmt"
	"fyne.io/fyne/v2"
	"strings"
	"time"
	"useful-tools/app/usefultools/view/constant"
)

type AesListKey struct {
	BaseKey
	Separator string
	CacheLen  int
}

func NewAesKeyList() *AesListKey {
	return &AesListKey{
		BaseKey: BaseKey{
			Name: constant.CacheKeyAesKeyList,
		},
		Separator: "!@!",
		CacheLen:  22,
	}
}

func (k *AesListKey) SetValue(name, key, iv string) (changed bool) {
	if fyne.CurrentApp().Preferences().Bool(constant.NavStatePreferenceSaveAesKey) {
		if name == "" {
			name = time.Now().String()
		}
		found := false
		content := k.formatContent(name, key, iv)
		keyList := fyne.CurrentApp().Preferences().StringList(k.Name)
		for i, ke := range keyList {
			keyInfo := strings.Split(ke, k.Separator)
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

		if len(keyList) > k.CacheLen {
			keyList = keyList[0:k.CacheLen]
		}

		fyne.CurrentApp().Preferences().SetStringList(k.Name, keyList)
	}
	return
}

func (k *AesListKey) GetValue() ([]string, map[string]string, map[string]string) {
	var names []string
	key := make(map[string]string)
	iv := make(map[string]string)
	keyList := fyne.CurrentApp().Preferences().StringList(k.Name)
	fmt.Println(keyList)
	for _, ke := range keyList {
		keyInfo := strings.Split(ke, k.Separator)
		if len(keyInfo) != 3 {
			continue
		}
		names = append(names, keyInfo[0])
		key[keyInfo[0]] = keyInfo[1]
		iv[keyInfo[0]] = keyInfo[2]
	}
	return names, key, iv
}

func (k *AesListKey) formatContent(values ...string) string {
	return strings.Join(values, k.Separator)
}

func (k *AesListKey) Clear() {
	fyne.CurrentApp().Preferences().SetStringList(k.Name, []string{})
}
