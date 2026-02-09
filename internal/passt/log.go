package passt

import (
	"log"
	"os"
)

// Logger corresponds to log.c, centralizing structured output.
type Logger struct {
	base *log.Logger
}

func NewLogger() *Logger {
	return &Logger{base: log.New(os.Stdout, "passt-go ", log.LstdFlags|log.Lmicroseconds)}
}

func (l *Logger) Infof(format string, args ...any) {
	l.base.Printf("INFO "+format, args...)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.base.Printf("ERROR "+format, args...)
}
