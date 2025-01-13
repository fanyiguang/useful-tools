package preferencekey

type BaseKey struct {
	Name string
}

func (b *BaseKey) Key() string {
	return b.Name
}
