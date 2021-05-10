package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

var instance *configData

type configData struct {
	address           string
	port              uint
	debug             bool
	chromedpWsAddress string
}

func init() {

	godotenv.Load()

	loggerPath := getEnv("LOGGER_PATH", "")
	if loggerPath != "" {
		if file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666); err == nil {
			log.SetOutput(file)
		} else {
			log.Println(fmt.Sprintf("couldn't open file with filename: %s, log output haven't been changed", loggerPath))
		}

	}

	instance = &configData{
		address:           getEnv("ADDRESS", "0.0.0.0"),
		port:              getEnvAsUInt("PORT", 1218),
		debug:             getEnvAsBool("DEBUG", false),
		chromedpWsAddress: getEnv("CHROMEDP_WS_ADDRESS", ""),
	}

	if instance.chromedpWsAddress != "" && !strings.HasPrefix(instance.chromedpWsAddress, "ws://") {
		instance.chromedpWsAddress = "ws://" + instance.chromedpWsAddress
	}

}

func Address() string {
	return instance.address
}

func Port() uint {
	return instance.port
}

func Debug() bool {
	return instance.debug
}

func ChromeDPWSAddress() string {
	return instance.chromedpWsAddress
}

func Version() string {
	return "1.0.0.1"
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsUInt(name string, defaultVal uint) uint {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return uint(value)
	}

	return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}
