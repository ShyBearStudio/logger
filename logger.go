package logger

type Logger interface {
	Traceln(args ...interface{})
	Tracef(format string, args ...interface{})
	Infoln(args ...interface{})
	Infof(format string, args ...interface{})
	Warningln(args ...interface{})
	Warningf(format string, args ...interface{})
	Errorln(args ...interface{})
	Errorf(format string, args ...interface{})
	Panicln(args ...interface{})
	Panicf(format string, args ...interface{})
}
