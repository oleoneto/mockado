package logger

import (
	"os"
	"strconv"
	"sync"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

var loggerOnce sync.Once

func GetLogger(name string) *logrus.Logger {
	loggerOnce.Do(func() {
		if logger != nil {
			return
		}

		logger = logrus.New()
		logger.SetLevel(func() logrus.Level {
			e := os.Getenv("LOG_LEVEL")
			i, err := strconv.Atoi(e)
			if err != nil {
				return logrus.InfoLevel
			}
			return logrus.Level(i)
		}())
	})

	return logger
}
