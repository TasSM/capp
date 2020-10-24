package util

import (
	"os"
)

func GetEnv(value string, def string) string {
	check := os.Getenv(value)
	if check != "" {
		return check
	}
	return def
}
