// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	Global Logger
    DEFAULT_LOGGER_NAME string = "Program"
)

func defaultconf() []byte {
	localDefaultConfig := []byte(`{
        "console": {
            "enable": true,
            "category": "Program",
            "level": "FINE",
            "pattern": "[%D %T] [%L] (%S) %M"
        },
        "files": [{
            "enable": true,
            "category": "Program",
            "level": "INFO",
            "filename":"../log/ProgramLog.log",
            "rotate": true,
            "maxsize": "800M",
            "daily": true,
            "pattern": "[%D %T] [%L] (%S) %M"
        }]
    }`)
	return localDefaultConfig
}

func init() {
	Global = NewDefaultLogger(FINE)
}

// Wrapper for (*Logger).LoadConfiguration
func LoadConfiguration(filename string, types ...string) {
	if len(types) > 0 && types[0] == "xml" {
		Global.LoadConfiguration(filename)
	} else {
		Global.LoadJsonConfiguration(filename)
	}
}

// Wrapper for (*Logger).AddFilter
func AddFilter(name string, lvl Level, writer LogWriter) {
	Global.AddFilter(name, lvl, writer)
}

func ChangeFilterLevel(name string, lvl string) {
    Global.ChangeFilterLevel(name, L4gGetLogLevel(lvl))
}

// Wrapper for (*Logger).Close (closes and removes all logwriters)
func Close() {
	Global.Close()
}

func Crash(args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(CRITICAL, strings.Repeat(" %v", len(args))[1:], args...)
	}
	panic(args)
}

// Logs the given message and crashes the program
func Crashf(format string, args ...interface{}) {
	Global.intLogf(CRITICAL, format, args...)
	Global.Close() // so that hopefully the messages get logged
	panic(fmt.Sprintf(format, args...))
}

// Compatibility with `log`
func Exit(args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(ERROR, strings.Repeat(" %v", len(args))[1:], args...)
	}
	Global.Close() // so that hopefully the messages get logged
	os.Exit(0)
}

// Compatibility with `log`
func Exitf(format string, args ...interface{}) {
	Global.intLogf(ERROR, format, args...)
	Global.Close() // so that hopefully the messages get logged
	os.Exit(0)
}

// Compatibility with `log`
func Stderr(args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(ERROR, strings.Repeat(" %v", len(args))[1:], args...)
	}
}

// Compatibility with `log`
func Stderrf(format string, args ...interface{}) {
	Global.intLogf(ERROR, format, args...)
}

// Compatibility with `log`
func Stdout(args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(INFO, strings.Repeat(" %v", len(args))[1:], args...)
	}
}

// Compatibility with `log`
func Stdoutf(format string, args ...interface{}) {
	Global.intLogf(INFO, format, args...)
}

// Send a log message manually
// Wrapper for (*Logger).Log
func Log(lvl Level, source, message string) {
	Global.Log(lvl, source, message)
}

// Send a formatted log message easily
// Wrapper for (*Logger).Logf
func Logf(lvl Level, format string, args ...interface{}) {
	Global.intLogf(lvl, format, args...)
}

// Send a closure log message
// Wrapper for (*Logger).Logc
func Logc(lvl Level, closure func() string) {
	Global.intLogc(lvl, closure)
}

// Utility for finest log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Finest
func Finest(arg0 interface{}, args ...interface{}) {
	const (
		lvl = FINEST
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Utility for fine log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Fine
func Fine(arg0 interface{}, args ...interface{}) {
	const (
		lvl = FINE
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Utility for debug log messages
// When given a string as the first argument, this behaves like Logf but with the DEBUG log level (e.g. the first argument is interpreted as a format for the latter arguments)
// When given a closure of type func()string, this logs the string returned by the closure iff it will be logged.  The closure runs at most one time.
// When given anything else, the log message will be each of the arguments formatted with %v and separated by spaces (ala Sprint).
// Wrapper for (*Logger).Debug
func Debug(arg0 interface{}, args ...interface{}) {
	const (
		lvl = DEBUG
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Utility for trace log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Trace
func Trace(arg0 interface{}, args ...interface{}) {
	const (
		lvl = TRACE
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Utility for info log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Info
func Info(arg0 interface{}, args ...interface{}) {
	const (
		lvl = INFO
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Utility for warn log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Warn
func Warn(arg0 interface{}, args ...interface{}) error {
	const (
		lvl = WARNING
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
		return errors.New(fmt.Sprintf(first, args...))
	case func() string:
		// Log the closure (no other arguments used)
		str := first()
		Global.intLogf(lvl, "%s", str)
		return errors.New(str)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
		return errors.New(fmt.Sprint(first) + fmt.Sprintf(strings.Repeat(" %v", len(args)), args...))
	}
	return nil
}

// Utility for error log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Error
func Error(arg0 interface{}, args ...interface{}) error {
	const (
		lvl = ERROR
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
		return errors.New(fmt.Sprintf(first, args...))
	case func() string:
		// Log the closure (no other arguments used)
		str := first()
		Global.intLogf(lvl, "%s", str)
		return errors.New(str)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
		return errors.New(fmt.Sprint(first) + fmt.Sprintf(strings.Repeat(" %v", len(args)), args...))
	}
	return nil
}

// Utility for critical log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Critical
func Critical(arg0 interface{}, args ...interface{}) error {
	const (
		lvl = CRITICAL
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
		return errors.New(fmt.Sprintf(first, args...))
	case func() string:
		// Log the closure (no other arguments used)
		str := first()
		Global.intLogf(lvl, "%s", str)
		return errors.New(str)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
		return errors.New(fmt.Sprint(first) + fmt.Sprintf(strings.Repeat(" %v", len(args)), args...))
	}
	return nil
}

func getAbsPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		fmt.Println("getAbsPath", err)
	}
	return abs
}

func GetParentDir(path string) string {
	lastUrl := filepath.Join(filepath.Dir(path), "..")
	return getAbsPath(lastUrl)
}

// prompath string is the program you running, abs path
func getValidLogPath(prompath string) (string, error) {
	fmt.Println("Program Path", getAbsPath(prompath))
	filename := filepath.Base(prompath) + ".log"
	// First Path Would Choose Parent Dir
	// >>
	//   >>bin
	//     >>prom
	//   >>log
	//     >>prom.log
	var tmpPath string = GetParentDir(prompath)
	firstChoose := filepath.Join(tmpPath, "log")

	// Second Path Would Choose Current Dir
	// --
	//   --prom
	//   --log
	//     --prom.log
	tmpPath = getAbsPath(filepath.Dir(prompath))
	secondChoose := filepath.Join(tmpPath, "log")

	backupDir := []string{firstChoose, secondChoose, "/data/log", getAbsPath(".")}
	fmt.Print("Dir Choose Pool[From Top to Down]:\n    ", strings.Join(backupDir, ",\n    "))
	fmt.Println("")
	for _, eachDir := range backupDir {
		// Check each Dir is Exists
		if _, err := os.Stat(eachDir); os.IsNotExist(err) {
			fmt.Println("log dir use fail :", err, "try another dir")
			continue
		}
		secondPath := filepath.Join(eachDir, filename)
		return secondPath, nil
	}
	return "", errors.New("no such path")
}

// If error exists, it must be bad thing happend
// level FINE, DEBUG, INFO, ...
func SetUniqueLogName(program string, level string) (string, error) {
	var localDefaultConfig LogConfig
	json.Unmarshal(defaultconf(), &localDefaultConfig)
	LogPath, err := getValidLogPath(program)
	fmt.Println("Select", LogPath, "as log dir")

	if err != nil {
		emsg := fmt.Sprint("No Valid Path Can Put The log: %s", err)
		fmt.Println(emsg)
		return "", errors.New(emsg)
	}

	localDefaultConfig.Files[0].Filename = LogPath
	localDefaultConfig.Files[0].Level = level
	localDefaultConfig.Console.Level = level
    if level != "DEBUG" && level != "FINE" {
        localDefaultConfig.Console.Enable = false
        fmt.Println("Log level is not FINE OR DEBUG, will not print log to the screen")
    }
	data, _ := json.Marshal(localDefaultConfig)
	var timestr string = strconv.Itoa(int(time.Now().Unix()))
	var tmpconfpath string = "/data/.tmpconf.json." + timestr
	err = ioutil.WriteFile(tmpconfpath, data, 0644)
	if err != nil {
		emsg := fmt.Sprint("Dump json config Fail :", err)
		fmt.Println(emsg)
		return "", errors.New(emsg)
	}
	LoadConfiguration(tmpconfpath)
	os.Remove(tmpconfpath)
	return program, nil
}
