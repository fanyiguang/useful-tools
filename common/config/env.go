package config

var env = "release"

func Env() string {
	return env
}

func IsTest() bool {
	return env == "test"
}
