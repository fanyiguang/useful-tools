package adapter

type PreferenceKey interface {
	Key() string
	Value() string
	SetValue(value string)
}
