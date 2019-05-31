package cache

import (
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	logger      *logrus.Entry
	loggerMutex sync.RWMutex
)

// SetLogger caches a logrus entry for future use
func SetLogger(s *logrus.Entry) {
	loggerMutex.Lock()
	logger = s
	loggerMutex.Unlock()
}

// GetLogger returns a cached logrus entry
func GetLogger() *logrus.Entry {
	loggerMutex.RLock()
	defer loggerMutex.RUnlock()

	return logger
}

// HasLogger returns true if we have a logger cached
func HasLogger() bool {
	loggerMutex.RLock()
	defer loggerMutex.RUnlock()

	return logger != nil
}
