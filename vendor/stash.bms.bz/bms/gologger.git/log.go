package gologger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

//================================================================================
// Struct
//================================================================================
type log struct {
	Host          string                 `json:"host"`
	App           string                 `json:"app"`
	File          string                 `json:"file"`
	Method        string                 `json:"method"`
	Line          int                    `json:"line"`
	LType         string                 `json:"type"`
	Component     string                 `json:"component"`
	Code          string                 `json:"code"`
	Description   string                 `json:"description"`
	Category      string                 `json:"category"`
	Doc           string                 `json:"doc"`
	TS            time.Time              `json:"ts"`
	Ref           map[string]interface{} `json:"ref"`
	DC            string                 `json:"dc"`
	Country       string                 `json:"country"`
	StackTrace    interface{}            `json:"stacktrace"`
	serializedLog string
}

// New Creates a log instance
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
func new(host string, app string, file string, method string, line int, lType string, component string, code string, description string, category string, doc string, ref map[string]interface{}) *log {
	stack := ref["stacktrace"]
	delete(ref, "stacktrace")
	return &log{
		Host:        host,
		App:         app,
		File:        file,
		Method:      method,
		Line:        line,
		LType:       lType,
		Component:   component,
		Code:        code,
		Description: description,
		Category:    category,
		Doc:         doc,
		TS:          time.Now().UTC(),
		DC:          os.Getenv("DC"),
		Country:     os.Getenv("COUNTRY"),
		Ref:         ref,
		StackTrace:  stack,
	}
}

// Serializes log
//
//InputParameters:
//
//Outputparameter:
//
func (lg *log) Serialize() (err error) {
	serializedLog, err := json.Marshal(lg)
	if err != nil {
		return err
	}

	lg.serializedLog = string(serializedLog)

	return nil
}

// Writes log
//
//InputParameters:
// sync: bool value specifies execution type(sync/async)
//
//Outputparameter:
//
func (lg *log) Write(sync bool) {
	if sync {
		fmt.Println(lg.serializedLog)
	} else {
		go fmt.Println(lg.serializedLog)
	}
}
