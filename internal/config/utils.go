package config

import (
	"log"
	"os"
	"strconv"
)

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("env %q value not found", key)
	}
	return value
}

func getEnvAsInt(key string) int {
	value := getEnv(key)
	v, e := strconv.Atoi(value)
	if e != nil {
		log.Fatalf("env %q value invalid int", key)
	}
	return v
}
