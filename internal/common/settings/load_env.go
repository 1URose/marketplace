package settings

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func GetEnvSrt(key string) string {
	v := os.Getenv(key)

	log.Printf("[postgresql] env %s: %s", key, v)

	if v == "" {
		log.Panicf("[postgresql][FATAL] env %s not set", key)
	}

	return v
}

func GetEnvInt32(key string) (int32, error) {
	valueStr := GetEnvSrt(key)

	value, err := strconv.ParseInt(valueStr, 10, 32)

	if err != nil {

		log.Panicf("[postgresql][FATAL] Invalid value for %s: %v", key, err)
	}

	return int32(value), nil
}

func GetEnvInt(key string) (int, error) {
	s := GetEnvSrt(key)
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid value for %s: %w", key, err)
	}
	return v, nil
}
