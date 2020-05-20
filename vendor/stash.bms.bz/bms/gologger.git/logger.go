package gologger

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
)

// LType logger supported enum type
type LType string

// Component logger supported enum type
type Component string

// Category logger supported enum type
type Category string

// Components category log as per logger
const (
	Fatal   LType = "fatal"
	Warning LType = "warning"
	Errsev1 LType = "errsev1"
	Errsev2 LType = "errsev2"
	Errsev3 LType = "errsev3"
	Info    LType = "info"
	Debug   LType = "debug"

	Datastore        Component = "Datastore"
	Queue            Component = "Queue"
	InternalServices Component = "InternalServices"
	ExternalServices Component = "ExternalServices"
	Application      Component = "Application"

	ConnectionError   Category = "ConnectionError"
	TimeoutError      Category = "TimeoutError"
	DataNotFoundError Category = "DataNotFoundError"
	ParseError        Category = "ParseError"
	ValidationError   Category = "ValidationError"
	UnknownError      Category = "UnknownError"
)

// Logger struct
type Logger struct {
	host  string
	app   string
	sync  bool
	trace struct {
		// Determain whether the stack trace need to be enabled or not
		Enable bool
		// Stack trace level
		Level int
	}
}

// New Creates a logger instance
//
//InputParameters:
//	 host: host of the app
//	 app: name of the app
//Outputparameter:
//	*logger: returns an instance of logger
func New(app string, sync ...bool) *Logger {
	// Get Host
	host, err := getHost()
	if err != nil {
		handleErr(errors.New("Fetching Host failed"))
	}

	syncVal := false
	if len(sync) == 1 {
		syncVal = sync[0]
	}

	// Instantiate Logger
	return &Logger{
		app:  app,
		host: host,
		sync: syncVal,
	}
}

// Log creates new logs instance
//
//InputParameters:
// file: name of the file
// method: name of the method
// line: line on which log occurs
// type: type of the log
// component: component to which the log belongs
// code: error code of the log
// description: description of the error
// category: category of the log
// doc: documentation link relevant to the log
// ref: reference object for app-level details
//Outputparameter:
//	status : returns status
func (lgr *Logger) Log(lType LType, component Component, code string, description string, category Category, doc string, ref map[string]interface{}, sync ...bool) {
	// Get Runtime Details (File, Line, Method)
	file, line, ok := getFileAndLine()
	if !ok {
		handleErr(errors.New("Fetching runtime details failed"))
	}

	method := getMethod()
	// If stack trace is enabled only then debug stacktrace will be added to console logs
	// Note:
	//		1) If Log type is fatal defult stack trace will added
	// 		2) If stack trace is enabled and log type is not debug,info,warning stack trace will added
	if (lgr.trace.Enable || lType == Fatal) &&
		!(lType == Debug || lType == Info || lType == Warning) {
		// (also called stack backtrace or stack traceback) is a report of the active stack frames
		// at a certain point in time during the execution of a program
		arr := getStackTrace(lgr.trace.Level)
		var stackTrace []stackFrame
		for _, val := range arr {
			stack := newStackFrame(val)
			// appending file path includes api and pkg here

			if !strings.Contains(stack.File, "/vendor/") {
				pkg := stack.Package
				if !(pkg == "runtime" || pkg == "goexit" || pkg == "build") {
					stackTrace = append(stackTrace, stack)
				}
			}
		}
		// To avoid panic
		if ref == nil {
			ref = make(map[string]interface{}, 0)
		}
		ref["stacktrace"] = stackTrace
	}
	// Log Instantiation
	newLog := new(lgr.host, lgr.app, file, method, line, string(lType), string(component), code, description, string(category), doc, ref)

	// Log operations
	err := newLog.Serialize()
	if err != nil {
		handleErr(err)
		return
	}

	if lType == Fatal {
		newLog.Write(true)
	} else {
		syncVal := lgr.sync
		if len(sync) == 1 {
			syncVal = sync[0]
		}
		newLog.Write(syncVal)
	}

	// if err != nil {
	// 	handleErr(err)
	// 	return
	// }
}

// StackTrace determines whether the stacktrace need to enabled or not, with stack trace level
func (lgr *Logger) StackTrace(enable bool, level ...int) {
	lgr.trace.Enable = enable
	lgr.trace.Level = 0
	if len(level) > 0 {
		lgr.trace.Level = level[0]
	}
}

// Get Host name
func getHost() (host string, err error) {
	return os.Hostname()
}

// Get File name and Line number
func getFileAndLine() (file string, line int, ok bool) {
	_, file, line, ok = runtime.Caller(2)
	if !ok {
		return "", 0, ok
	}

	return file, line, ok
}

// GetMethod
func getMethod() (method string) {
	// we get the callers as interprets - but we just need 1
	fPcs := make([]uintptr, 1)

	// skip 3 levels to get to the caller of whoever called Caller()
	n := runtime.Callers(3, fPcs)
	if n == 0 {
		return "n/a" // proper error her would be better
	}

	// get the info of the actual function that's in the pointer
	methodObj := runtime.FuncForPC(fPcs[0] - 1)
	if methodObj == nil {
		return "n/a"
	}

	// return its name
	return methodObj.Name()
}

// Error handling for lib
func handleErr(err error) {
	fmt.Println("An error occurred in log lib:", err)
}
