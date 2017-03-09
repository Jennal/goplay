package log

import (
	"fmt"
	l "log"
	"os"
)

type Logger interface {
	Log(line string)
	Logf(format string, args ...interface{})
	Error(err error)
	Errorf(format string, args ...interface{})
}

func setStdout() {
	l.SetFlags(l.LstdFlags | l.Lshortfile)
	l.SetOutput(os.Stdout)
}

func setStderr() {
	l.SetFlags(l.LstdFlags | l.Lshortfile)
	l.SetOutput(os.Stderr)
}

type _logger struct {
	prefix string
	depth  int
}

func NewLogger(prefix string, depth int) Logger {
	if len(prefix) > 0 {
		prefix = prefix + " "
	}

	return _logger{
		prefix: prefix,
		depth:  depth,
	}
}

func (logger _logger) Log(line string) {
	setStdout()
	l.Output(logger.depth, logger.prefix+line)
}

func (logger _logger) Logf(format string, args ...interface{}) {
	setStdout()
	l.Output(logger.depth, logger.prefix+fmt.Sprintf(format, args...))
}

func (logger _logger) Error(err error) {
	setStderr()
	l.Output(logger.depth, logger.prefix+err.Error())
}

func (logger _logger) Errorf(format string, args ...interface{}) {
	setStderr()
	l.Output(logger.depth, logger.prefix+fmt.Sprintf(format, args...))
}
