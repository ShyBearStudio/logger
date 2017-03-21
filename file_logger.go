package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	TraceFileName   string = "trace.log"
	InfoFileName    string = "info.log"
	WarningFileName string = "warning.log"
	ErrorFileName   string = "error.log"
	LogFormat       int    = log.Ldate | log.Ltime | log.Lshortfile
)

var logDescs = []struct {
	fileName string
	prefix   string
}{
	{TraceFileName, "TRACE: "},
	{InfoFileName, "INFO: "},
	{WarningFileName, "WARN: "},
	{ErrorFileName, "ERROR: "},
}

type FileLogger struct {
	dir   string
	log   map[string]*log.Logger
	files map[string]*os.File
}

func NewFileLogger(logDir string) (*FileLogger, error) {
	var fileLogger = new(FileLogger)
	fileLogger.log = make(map[string]*log.Logger)
	fileLogger.files = make(map[string]*os.File)
	for _, desc := range logDescs {
		file, err := openLogFile(logDir, desc.fileName)
		if err != nil {
			return nil, err
		}
		fileLogger.files[desc.fileName] = file
		fileLogger.log[desc.fileName] = log.New(file, desc.prefix, LogFormat)
	}
	return fileLogger, nil
}

func (fl *FileLogger) Release() {
	for _, desc := range logDescs {
		file := fl.files[desc.fileName]
		file.Close()
	}
}

func openLogFile(dir, name string) (file *os.File, err error) {
	os.MkdirAll(dir, os.ModePerm)
	path := filepath.Join(dir, name)
	file, err = os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	return
}

func (fl *FileLogger) println(fileName string, args ...interface{}) {
	l := fl.log[fileName]
	l.Output(3, fmt.Sprintln(args...))
}

func (fl *FileLogger) printf(fileName string, format string, args ...interface{}) {
	l := fl.log[fileName]
	l.Output(3, fmt.Sprintf(format, args...))
}

func (fl *FileLogger) Traceln(args ...interface{}) {
	fl.println(TraceFileName, args...)
}

func (fl *FileLogger) Tracef(format string, args ...interface{}) {
	fl.printf(TraceFileName, format, args...)
}

func (fl *FileLogger) Infoln(args ...interface{}) {
	fl.println(InfoFileName, args...)
}

func (fl *FileLogger) Infof(format string, args ...interface{}) {
	fl.printf(InfoFileName, format, args...)
}

func (fl *FileLogger) Warningln(args ...interface{}) {
	fl.println(WarningFileName, args...)
}

func (fl *FileLogger) Warningf(format string, args ...interface{}) {
	fl.printf(WarningFileName, format, args...)
}
func (fl *FileLogger) Errorln(args ...interface{}) {
	fl.println(ErrorFileName, args...)
}

func (fl *FileLogger) Errorf(format string, args ...interface{}) {
	fl.printf(ErrorFileName, format, args...)
}

func (fl *FileLogger) Panicln(args ...interface{}) {
	fl.println(ErrorFileName, args...)
	panic(fmt.Sprintln(args...))
}

func (fl *FileLogger) Panicf(format string, args ...interface{}) {
	fl.printf(ErrorFileName, format, args...)
	panic(fmt.Sprintf(format, args...))
}
