package main

import (
	"bufio"
	"fmt"
	"os"
)

// Logger logs
type Logger struct {
	file   *os.File
	writer *bufio.Writer
}

// NewLogger returns a new Logger. File is a path to file on the disk where
// logs should be saved. It does not have to exist already. If it doesn't,
// it will be created.
func NewLogger(filepath string) *Logger {
	var logger Logger

	// errors are truncated which means that if something goes wrong (ie no
	// permissions to create the file), it will silently fail
	logger.file, _ = os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0755)
	logger.file.Truncate(0)

	logger.writer = bufio.NewWriter(logger.file)
	return &logger
}

func (l *Logger) Clear(text string) {
	l.writer.Flush()
}

func (l *Logger) formatStr(evt EventType, text string) string {
	return fmt.Sprintf("[%s] %s\n", evt, text)
}

func (l *Logger) Log(evt EventType, text string) {
	l.writer.WriteString(l.formatStr(evt, text))
}

func (l *Logger) Close() {
	// flush the buffer aka write anything left
	// in the buffer into the file
	l.writer.Flush()
	// close the file. Error checks are not performed.
	l.file.Close()
}
