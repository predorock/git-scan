package logstreamer

import (
	"log"
	"os"
	"sync"
)

var logger *Logger
var once sync.Once

type Logger struct {
	Log    *log.Logger
	Stderr Logstreamer
	Stdout Logstreamer
}

func GetInstance() *Logger {
	once.Do(func() {
		logger = createLogger()
	})
	return logger
}

func createLogger() *Logger {

	logger := log.New(os.Stdout, "--> ", log.Ldate|log.Ltime)

	// Setup a streamer that we'll pipe cmd.Stdout to
	logStreamerOut := NewLogstreamer(logger, "stdout", "scan", false)
	defer logStreamerOut.Close()
	// Setup a streamer that we'll pipe cmd.Stderr to.
	// We want to record/buffer anything that's written to this (3rd argument true)
	logStreamerErr := NewLogstreamer(logger, "stderr", "scan", false)
	defer logStreamerErr.Close()

	return &Logger{
		Log:    logger,
		Stderr: *logStreamerErr,
		Stdout: *logStreamerOut,
	}
}
