package tlog

import (
	"fmt"
	"log"
	"os"
)

var Fatal func(args ...interface{})
var Fatalf func(format string, args ...interface{})

var Error func(args ...interface{})
var Errorf func(format string, args ...interface{})

var Warning func(args ...interface{})
var Warningf func(format string, args ...interface{})

var Info func(args ...interface{})
var Infof func(format string, args ...interface{})

var Trace func(args ...interface{})
var Tracef func(format string, args ...interface{})

var Debug func(args ...interface{})
var Debugf func(format string, args ...interface{})

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	setDefLogs()
}

func setDefLogs() {
	Fatal, Fatalf = levelLogFuncs("[fatal] ")
	Error, Errorf = levelLogFuncs("[error] ")
	Warning, Warningf = levelLogFuncs("[warning] ")
	Info, Infof = levelLogFuncs("[info] ")
	Trace, Tracef = levelLogFuncs("[trace] ")
	Debug, Debugf = levelLogFuncs("[debug] ")
}

func levelLogFuncs(prefix string) (func(...interface{}), func(string, ...interface{})) {
	logger := log.New(os.Stderr, prefix, log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	log := func(a ...interface{}) {
		logger.Output(2, fmt.Sprint(a...))
	}
	logf := func(format string, a ...interface{}) {
		logger.Output(2, fmt.Sprintf(format, a...))
	}
	return log, logf
}

func NoLog(args ...interface{}) {
}

func NoLogf(format string, args ...interface{}) {
}
