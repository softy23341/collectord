// Licensed under the LGPLv3, see LICENCE file for details.

package erry

import (
	"fmt"
	"runtime"
)

// Category TBD
type Category struct {
	parent *Category
}

// GetCategory TBD
func GetCategory(e error) *Category {
	if categoried, ok := e.(Categorier); ok {
		return categoried.Category()
	}
	return nil
}

// GetTopCategory TBD
func GetTopCategory(e error) *Category {
	c := GetCategory(e)
	for check := c; check != nil; check = check.parent {
		if check.IsTop() {
			return check
		}
	}
	return nil
}

// NewCategory TBD
func NewCategory() *Category {
	return &Category{}
}

// IsTop TBD
func (c *Category) IsTop() bool {
	return c.parent == nil
}

// IsSub TBD
func (c *Category) IsSub() bool {
	return !c.IsTop()
}

// New TBD
func (c *Category) New(msg string) *Er {
	er := &Er{msg: msg, cat: c}
	return er
}

// Newf TBD
func (c *Category) Newf(format string, args ...interface{}) *Er {
	er := &Er{msg: fmt.Sprintf(format, args...), cat: c}
	return er
}

// NewSubCategory TBD
func (c *Category) NewSubCategory() *Category {
	return &Category{parent: c}
}

// Grab TBD
func (c *Category) Grab(e error) *Er {
	return &Er{cat: c, msg: e.Error()}
}

// Wrap TBD
func (c *Category) Wrap(e error) *Er {
	return &Er{cat: c, cause: e, msg: e.Error()}
}

// WrapRaw TBD
func (c *Category) WrapRaw(e error) *Er {
	if err, ok := e.(*Er); ok {
		return err
	}
	return c.Wrap(e)
}

// Contains TBD
func (c *Category) Contains(e error) bool {
	cat := GetCategory(e)
	for check := cat; check != nil; check = cat.parent {
		if check == c {
			return true
		}
	}
	return false
}

// Categorier TBD
type Categorier interface {
	error
	Category() *Category
}

// Causer TBD
type Causer interface {
	Cause() error
}

// Er TBD
type Er struct {
	cat   *Category
	msg   string
	cause error

	// file and line hold the source code location where the error was
	// created.
	file string
	line int
}

// Location is the file and line of where the error was most recently
// created or annotated.
func (e *Er) Location() (filename string, line int) {
	return e.file, e.line
}

func (e *Er) Error() string {
	return e.msg
}

// Category TBD
func (e *Er) Category() *Category {
	return e.cat
}

// Cause TBD
func (e *Er) Cause() error {
	return e.cause
}

// Message returns the message stored with the most recent location. This is
// the empty string if the most recent call was Trace, or the message stored
// with Annotate or Mask.
func (e *Er) Message() string {
	return e.msg
}

// Trace TBD
func (e *Er) Trace() *Er {
	*e = *trace(e.dup())
	e.setLocation(1)
	return e
}

// Annotate TBD
func (e *Er) Annotate(message string) *Er {
	*e = *annotate(e.dup(), message)
	e.setLocation(1)
	return e

}

// Annotatef TBD
func (e *Er) Annotatef(format string, args ...interface{}) *Er {
	*e = *annotatef(e.dup(), format, args)
	e.setLocation(1)
	return e
}

func (e *Er) dup() *Er {
	clone := *e
	return &clone
}

// SetLocation records the source location of the error at callDepth stack
// frames above the call.
func (e *Er) setLocation(callDepth int) {
	_, file, line, _ := runtime.Caller(callDepth + 1)
	e.file = trimGoPath(file)
	e.line = line
}

// StackTrace returns one string for each location recorded in the stack of
// errors. The first value is the originating error, with a line for each
// other annotation or tracing of the error.
func (e *Er) StackTrace() []string {
	return errorStack(e)
}

// GetCause TBD
func GetCause(e error) error {
	if causer, ok := e.(Causer); ok {
		return causer.Cause()
	}
	return nil
}
