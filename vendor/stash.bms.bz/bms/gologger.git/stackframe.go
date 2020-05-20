package gologger

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
)

type (
	// stackTrace is a stack trace.
	stackTrace []uintptr

	// A StackFrame contains all necessary information about to generate a line
	// in a callstack.
	stackFrame struct {
		// The path to the file containing this ProgramCounter
		File string `json:"file"`
		// The LineNumber in that file
		LineNumber int `json:"linenumber"`
		// The Name of the function that contains this ProgramCounter
		Name string `json:"name"`
		// The Package that contains this function
		Package string `json:"package"`
		// The underlying ProgramCounter
		ProgramCounter uintptr `json:"-"`
	}
)

// newStackFrame populates a stack frame object from the program counter.
func newStackFrame(pc uintptr) (frame stackFrame) {
	frame = stackFrame{ProgramCounter: pc}
	if frame.Func() == nil {
		return
	}
	frame.Package, frame.Name = packageAndName(frame.Func())

	// pc -1 because the program counters we use are usually return addresses,
	// and we want to show the line that corresponds to the function call
	frame.File, frame.LineNumber = frame.Func().FileLine(pc - 1)
	return
}

// Func returns the function that contained this frame.
func (frame *stackFrame) Func() *runtime.Func {
	if frame.ProgramCounter == 0 {
		return nil
	}
	return runtime.FuncForPC(frame.ProgramCounter)
}

// String returns the stackframe formatted in the same way as go does
// in runtime/debug.Stack()
func (frame *stackFrame) String() string {
	str := fmt.Sprintf("%s:%d (0x%x)\n", frame.File, frame.LineNumber, frame.ProgramCounter)

	source, err := frame.SourceLine()
	if err != nil {
		return str
	}

	return str + fmt.Sprintf("\t%s: %s\n", frame.Name, source)
}

// SourceLine gets the line of code (from File and Line) of the original source if possible.
func (frame *stackFrame) SourceLine() (string, error) {
	data, err := ioutil.ReadFile(frame.File)

	if err != nil {
		return "", err
	}

	lines := bytes.Split(data, []byte{'\n'})
	if frame.LineNumber <= 0 || frame.LineNumber >= len(lines) {
		return "???", nil
	}
	// -1 because line-numbers are 1 based, but our array is 0 based
	return string(bytes.Trim(lines[frame.LineNumber-1], " \t")), nil
}

func packageAndName(fn *runtime.Func) (string, string) {
	name := fn.Name()
	pkg := ""

	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//  runtime/debug.*T·ptrmethod
	// and want
	//  *T.ptrmethod
	// Since the package path might contains dots (e.g. code.google.com/...),
	// we first remove the path prefix if there is one.
	if lastslash := strings.LastIndex(name, "/"); lastslash >= 0 {
		pkg += name[:lastslash] + "/"
		name = name[lastslash+1:]
	}
	if period := strings.Index(name, "."); period >= 0 {
		pkg += name[:period]
		name = name[period+1:]
	}

	name = strings.Replace(name, "·", ".", -1)
	return pkg, name
}

// getStackTrace returns a new StackTrace.
func getStackTrace(skipFrames int) stackTrace {
	skip := 2 // skips runtime.Callers and this function
	skip += skipFrames
	callers := make([]uintptr, 100)
	written := runtime.Callers(skip, callers)
	return stackTrace(callers[0:written])
}
