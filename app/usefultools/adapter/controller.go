package adapter

type Controller interface {
	FormatJson(data string) string
	DoubleClicked() (res bool)
	ClearCache()
}
