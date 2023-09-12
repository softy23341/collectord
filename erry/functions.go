// Licensed under the LGPLv3, see LICENCE file for details.

package erry

import (
	"fmt"
	"strings"
)

// Trace adds the location of the Trace call to the stack.  The Cause of the
// resulting error is the same as the error parameter.  If the other error is
// nil, the result will be nil.
//
// For example:
//
//	if err := SomeFunc(); err != nil {
//	    return errors.Trace(err)
//	}
func Trace(other error) *Er {
	err := trace(other)
	if err != nil {
		err.setLocation(1)
	}
	return err
}

func trace(other error) *Er {
	if other == nil {
		return nil
	}
	category := GetCategory(other)
	err := &Er{cause: other, cat: category, msg: other.Error()}
	return err
}

// Annotate is used to add extra context to an existing error. The location of
// the Annotate call is recorded with the annotations. The file, line and
// function are also recorded.
//
// For example:
//
//	if err := SomeFunc(); err != nil {
//	    return errors.Annotate(err, "failed to frombulate")
//	}
func Annotate(other error, message string) *Er {
	err := annotate(other, message)
	if err != nil {
		err.setLocation(1)
	}
	return err
}

func annotate(other error, message string) *Er {
	if other == nil {
		return nil
	}
	category := GetCategory(other)
	err := &Er{
		cat:   category,
		cause: other,
		msg:   message,
	}
	return err
}

// Annotatef is used to add extra context to an existing error. The location of
// the Annotate call is recorded with the annotations. The file, line and
// function are also recorded.
//
// For example:
//
//	if err := SomeFunc(); err != nil {
//	    return errors.Annotatef(err, "failed to frombulate the %s", arg)
//	}
func Annotatef(other error, format string, args ...interface{}) *Er {
	err := annotatef(other, format, args)
	if err != nil {
		err.setLocation(1)
	}
	return err
}

func annotatef(other error, format string, args ...interface{}) *Er {
	if other == nil {
		return nil
	}
	category := GetCategory(other)
	err := &Er{
		cat:   category,
		cause: other,
		msg:   fmt.Sprintf(format, args...),
	}
	return err
}

// DeferredAnnotatef annotates the given error (when it is not nil) with the given
// format string and arguments (like fmt.Sprintf). If *err is nil, DeferredAnnotatef
// does nothing. This method is used in a defer statement in order to annotate any
// resulting error with the same message.
//
// For example:
//
//	defer DeferredAnnotatef(&err, "failed to frombulate the %s", arg)
func DeferredAnnotatef(err *error, format string, args ...interface{}) {
	if *err == nil {
		return
	}
	category := GetCategory(*err)
	newErr := &Er{
		cat:   category,
		msg:   fmt.Sprintf(format, args...),
		cause: *err,
	}
	newErr.setLocation(1)
	*err = newErr
}

type wrapper interface {
	// Message returns the top level error message,
	// not including the message from the Previous
	// error.
	Message() string

	Causer
}

type locationer interface {
	Location() (string, int)
}

var (
	_ wrapper    = (*Er)(nil)
	_ locationer = (*Er)(nil)
)

// ErrorStack returns a string representation of the annotated error. If the
// error passed as the parameter is not an annotated error, the result is
// simply the result of the Error() method on that error.
//
// If the error is an annotated error, a multi-line string is returned where
// each line represents one entry in the annotation stack. The full filename
// from the call stack is used in the output.
//
//	first error
//	github.com/juju/errors/annotation_test.go:193: first error
//	github.com/juju/errors/annotation_test.go:194: annotation
//	github.com/juju/errors/annotation_test.go:195: annotation
//	github.com/juju/errors/annotation_test.go:196: more context
//	github.com/juju/errors/annotation_test.go:197: more context
func ErrorStack(err error) string {
	return strings.Join(errorStack(err), "\n")
}

func errorStack(err error) []string {
	if err == nil {
		return nil
	}

	var nextError error

	// We want the first error first
	var lines []string
	for {
		nextError = nil
		if cerr, ok := err.(wrapper); ok {
			nextError = cerr.Cause()
		}

		var buff []byte
		pointed := false
		if err, ok := err.(locationer); ok {
			file, line := err.Location()

			file = trimGoPath(file)
			if file != "" {
				pointed = true

				buff = append(buff, fmt.Sprintf("%s:%d", file, line)...)
				buff = append(buff, ": "...)
			}
		}

		if cerr, ok := err.(wrapper); ok {
			if pointed || cerr.Cause() == nil {
				buff = append(buff, cerr.Message()...)
			}
		} else {
			buff = append(buff, err.Error()...)
		}

		if len(buff) != 0 {
			lines = append(lines, string(buff))
		}

		if err = nextError; err == nil {
			break
		}
	}

	// reverse the lines to get the original error, which was at the end of
	// the list, back to the start.
	var result []string
	for i := len(lines); i > 0; i-- {
		result = append(result, lines[i-1])
	}
	return result
}
