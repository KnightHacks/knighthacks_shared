package utils

import (
	"log"
	"os"
)

func GetEnvOrDie(key string) string {
	env, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("You must provide the %s environmental variable\n", key)
	}
	return env
}
