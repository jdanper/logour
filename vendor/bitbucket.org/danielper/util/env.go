package util

import "os"

//GetEnvOrDefault returns the environment value or a default value
func GetEnvOrDefault(key, def string) string {
	value := os.Getenv(key)
	if value == "" {
		return def
	}

	return value
}
