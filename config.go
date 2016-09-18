package main

import (
	"os"
	"strconv"
)

var (
	EnvDBAddress = GetEnv("MONGO_PORT_27017_TCP", "tcp://localhost:27017")
	EnvDBName    = GetEnv("DB_NAME", "gateway")
	EnvIPAddress = GetEnv("IP_ADDR", "0.0.0.0")
	EnvPort      = GetEnv("PORT", "80")
)

// GetEnv get an environment key as string, returns the default if not found
func GetEnv(key string, defaultVal string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultVal
	}
	return value
}

// GetEnvInt get an environment key as int, returns the default if not found
func GetEnvInt(key string, defaultVal int) int {
	valueStr := GetEnv(key, strconv.Itoa(defaultVal))
	res, err := strconv.Atoi(valueStr)
	if err != nil {
		res = defaultVal
	}
	return res
}

// GetEnvBool get an environment key as bool. If any value returns true,
// returns the default if not found
func GetEnvBool(key string, defaultVal bool) bool {
	valueStr := GetEnv(key, "")
	if valueStr != "" {
		return true
	}
	return false
}
