# Go Logger

## Installation

To install run the following command inside your project directory,

```bash
go get ssh://git@stash.bms.bz/bms/gologger.git
```

## Vendoring to your project (recommended for all projects)

```
$ cd <project root directory>/vendor/bms/
$ git clone ssh://git@stash.bms.bz/bms/gologger.git
```

Or use a Go vendoring tool to vendor.

E.g. [Dep](https://github.com/golang/dep), [GVT](https://github.com/FiloSottile/gvt)

## Usage Example

### Logger Initialization

```go
import lgr "stash.bms.bz/bms/gologger"
logger := lgr.New("TransactionEngine",true)
```

### Logging

```go
logger.Log("fatal", "datastore", "ab-101", "error fetching abc from xyz", "connectionError", "https://abc.xyz/", map[string]interface{}{"foo": "1", "bar": 2},true)
```

## Functions

<html><h3><u>logger package</u></h3></html>

### Logger Struct
	type Logger struct {
		host string
		app  string
		sync bool
	}

#### Fields :
Param        |Type   | Description
-------------|-------|----------------
app          |string |Name of the app
host         |string |host address of the app
sync         |bool   |Default flag to specify the execution type

<html><br/></html>

### New()
	func New(app string,sync ...bool) (logger *logger)

Creates a logger instance

#### Input Parameters :
Param        |Type    | Description
-------------|--------|---------------
app          |string  |Name of the app
sync         |...bool |Default flag to specify the execution type(Optional field)

#### Output :
Param        |Type    | Description
-------------|--------|---------------
logger       |*logger |Logger instance

<html><br/></html>

### Log()
	func (lgr *logger) Log(lType lg.LType, component lg.Component, code string, description string, category lg.Category, doc string, ref map[string]interface{}, sync ...bool) (err error)

Creates a logger instance

#### Input Parameters :
Param        |Type                  | Description
-------------|----------------------|---------------------------------------
type         |log.Type              |Type of the log
component    |log.Component         |Component to which the log belongs
code         |string                |Error code of the log
description  |string                |Decription of the error
category     |log.Category          |Category of the log
doc          |string                |Documentation link relevant to the log
ref          |map[string]interface{}|Reference object for app-level details
sync         |...bool               |Sync flag to specify the execution type(Optional field).Default value will assigned from New().

#### Output :
Param        |Type  | Description
-------------|------|--------------
err          |error |Error if any

<html><br/></html>

<html><h3><u>log package</u></h3></html>

### New()
	func New(file string, method string, line int, lType LType, component Component, code string, description string, category Category, doc string, ref map[string]interface{}) (log *log)

Creates a log instance

##### Input parameters:
Param        |Type                  | Description
-------------|----------------------|---------------------------------------
host         |string                |HostName of the server instance
app          |string                |Name of the application
file         |string                |Name of the file
method       |string                |Name of the method
line         |string                |Line on which log occurs
type         |log.Type              |Type of the log
component    |log.Component         |Component to which the log belongs
code         |string                |Error code of the log
description  |string                |Description of the error
category     |log.Category          |Category of the log
doc          |string                |Documentation link relevant to the log
ref          |map[string]interface{}|Reference object for app-level details

#### Output :
Param        |Type    | Description
-------------|--------|-------------
logger       |*log    |Log instance

<html><br/></html>

### serialize()
	func (lg *log) Serialize() (err error)

Serializes log

##### Input parameters:
Param        |Type                  | Description
-------------|----------------------|--------------
lg           |string                |log instance

<html><br/></html>

### write()
	func (lg *log) Write(sync bool) (err error)

Writes the log

##### Input parameters:
Param        |Type                  | Description
-------------|----------------------|--------------
lg           |string                |log instance
sync         |bool                  |bool value specifies execution type(sync/async)

#### Output :
Param        |Type    | Description
-------------|--------|-------------
logger       |*log    |Log instance
