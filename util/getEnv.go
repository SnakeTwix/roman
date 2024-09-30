package util

import (
	"log"
	"os"
)

func GetEnv(param string) string {
	value := os.Getenv(param)

	if value == "" {
		log.Fatalf("%v not specified \r\n", param)
	}

	return value
}
