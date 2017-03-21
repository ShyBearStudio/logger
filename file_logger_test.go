package logger

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testLogDir string = "TestLogDir"

func TestMain(m *testing.M) {
	if err := os.RemoveAll(testLogDir); err != nil {
		panic(fmt.Sprintf("Cannot clean up test log directory: '%v'", err))
	}
	code := m.Run()
	if err := os.RemoveAll(testLogDir); err != nil {
		panic(fmt.Sprintf("Cannot clean up test log directory: '%v'", err))
	}
	os.Exit(code)
}

func TestNewFileLogger(t *testing.T) {
	fl, err := NewFileLogger(testLogDir)
	if err != nil {
		t.Fatalf("Cannot create file logger: '%v'", err)
	}
	if _, err := os.Stat(testLogDir); os.IsNotExist(err) {
		t.Fatalf("Log directory '%s' does not exist", testLogDir)
	}
	for _, desc := range logDescs {
		fullName := fullFileName(desc.fileName)
		if _, err := os.Stat(fullName); os.IsNotExist(err) {
			t.Fatalf("Log file '%s' does not exist", fullName)
		}
	}
	assertLogger(t, fl)
	fl.Release()
}

func fullFileName(fileName string) string {
	return filepath.Join(testLogDir, fileName)
}

func assertLogger(t *testing.T, l Logger) {
	format := "val1 = %s --- val2 = %s"
	val1 := "val1"
	val2 := "val2"
	msg := fmt.Sprintf(format, val1, val2)
	l.Traceln(msg)
	l.Tracef(format, val1, val2)
	l.Infoln(msg)
	l.Infof(format, val1, val2)
	l.Warningln(msg)
	l.Warningf(format, val1, val2)
	l.Errorln(msg)
	l.Errorf(format, val1, val2)
	assertEachLineContains(t, TraceFileName, msg)
	assertEachLineContains(t, InfoFileName, msg)
	assertEachLineContains(t, WarningFileName, msg)
	assertEachLineContains(t, ErrorFileName, msg)
}

func assertEachLineContains(t *testing.T, fileName, msg string) {
	fullName := fullFileName(fileName)
	file, err := os.Open(fullName)
	if err != nil {
		t.Fatalf("Cannot open file '%s': %v", fullName, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, msg) {
			t.Fatalf("In log file '%s' the line is unexpected:\n\tline: '%s'\n\texpected sub-string: '%s'", fullName, line, msg)
		}
	}
}
