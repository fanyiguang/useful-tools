package str

import (
	"fmt"
	"testing"
)

func TestTrimInvalidCharacter(t *testing.T) {
	character := TrimInvalidCharacter("\n     dlfjakdjflka\n \t \r\n")
	fmt.Printf("%s\n", character)
}
