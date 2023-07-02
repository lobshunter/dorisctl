package log

import (
	"log"
)

func Debugf(format string, args ...any) {
	log.Printf("[DEBUG] "+format, args...)
}

func Infof(format string, args ...any) {
	log.Printf("[INFO] "+format, args...)
}

func Warnf(format string, args ...any) {
	log.Printf("[WARN] "+format, args...)
}

func Errorf(format string, args ...any) {
	log.Printf("[ERROR] "+format, args...)
}

func Fatalf(format string, args ...any) {
	log.Fatalf("[FATAL] "+format, args...)
}
