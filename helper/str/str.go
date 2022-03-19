package str

import "strings"

func TrimStringSpace(list ...string) (newList []string) {
	for _, str := range list {
		newList = append(newList, strings.TrimSpace(str))
	}
	return
}

func TrimInvalidCharacter(str string) string {
	return strings.TrimSpace(str)
}
