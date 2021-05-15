package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var instance *configData

type mongoDBConfig struct {
	Hostname string
	Port     int
	Username string
	Password string
}

type configData struct {
	mongoDBConfig     *mongoDBConfig
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

	useMongoDB := getEnvAsBool("MONGODB_USE", false)
	if useMongoDB {
		mc := mongoDBConfig{
			Hostname: getEnv("MONGODB_HOSTNAME", ""),
			Port:     int(getEnvAsUInt("MONGODB_PORT", 0)),
			Username: getEnv("MONGODB_USERNAME", ""),
			Password: getEnv("MONGODB_PASSWORD", ""),
		}
		instance.mongoDBConfig = &mc
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

func MongoDBConfig() *mongoDBConfig {
	return instance.mongoDBConfig
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
