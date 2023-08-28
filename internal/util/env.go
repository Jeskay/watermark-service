package util

import "os"

func EnvString(env, default_value string) string {
	e := os.Getenv(env)
	if e == "" {
		return default_value
	}
	return e
}
