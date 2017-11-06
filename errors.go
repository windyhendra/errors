// Errors package is to defined custom error golang package
// This package will add some contexts or informations for the error
// The purpose is to make developer easier when read the errors
// Also there're http error response, so it can defines what will return as http response when the system found the error
package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// Errs struct is to defined standard errors
type Errs struct {
	err  error
	file string
	line int

	// errHTTP is defined to store data for http response purpose
	errHTTP *ErrsHTTP

	// fields is defined to store additional error context data
	fields Fields
}

// Fields will store data as map
type Fields map[string]interface{}

// ErrsHTTP struct will be used as http error response type
type ErrsHTTP struct {
	statusCode int
	message    string
}

// fileDelim will be used as stopword for caller runtime
const fileDelim = "errors/"

// New function will be used as initialized function
// args can be string, *Errs, error, Fields, *ErrsHTTP data type
// another data type can not be used for this error
// mandatory data type that should be there are string or *Errs or error
// those data type will be used to define error, so the system will know that there's error
// ErrsHTTP and Fields only used for additional context data purpose
func New(args ...interface{}) *Errs {
	var errs Errs

	for _, arg := range args {
		switch arg.(type) {
		case string:
			errs.err = errors.New(arg.(string))
		case *Errs:
			errs = *arg.(*Errs)
		case error:
			errs.err = arg.(error)
		case Fields:
			errs.fields = arg.(Fields)
		case *ErrsHTTP:
			errs.errHTTP = arg.(*ErrsHTTP)
		default:
			// unknown error
			_, file, line, ok := runtime.Caller(1)

			if ok {
				slash := strings.LastIndex(file, fileDelim)
				file = file[slash+len(fileDelim):]
			}

			fmt.Printf("errors.Errs: bad call from %s:%d: %v", file, line, args)
		}
	}

	if !errs.isNil() {
		var ok bool
		_, errs.file, errs.line, ok = runtime.Caller(1)

		if ok {
			slash := strings.LastIndex(errs.file, fileDelim)
			errs.file = errs.file[slash+len(fileDelim):]
		}

		return &errs
	}

	return nil
}

// isNil will return nil if there's no error found
func (err Errs) isNil() bool {
	if err.err != nil {
		return false
	}
	return true
}

// Error function will return error string
func (err Errs) Error() string {
	var result string

	if err.GetCaller() != "" {
		result = fmt.Sprintf("source=%s", err.GetCaller())
	}

	if err.GetMessage() != "" {
		result = fmt.Sprintf("%s\tmessage=%s", result, err.GetMessage())
	}

	if err.GetFieldsString() != "" {
		result = fmt.Sprintf("%s\t%s", result, err.GetFieldsString())
	}

	return result
}

// GetMessage will return error message
func (err Errs) GetMessage() string {
	var message string
	if err.err != nil {
		message = err.err.Error()
	}
	return message
}

// GetCaller will return caller runtime string
// using format 'file:line'
func (err Errs) GetCaller() string {
	var caller string
	if err.file != "" && err.line != 0 {
		caller = fmt.Sprintf("%s:%d", err.file, err.line)
	}
	return caller
}

// WithHTTPError is to set ErrsHTTP data
func WithHTTPError(httpCode int, message string) *ErrsHTTP {
	return &ErrsHTTP{statusCode: httpCode, message: message}
}

// GetFieldsString return fields that already convert into string
func (err Errs) GetFieldsString() string {
	fields := ""
	if err.fields != nil {
		for key, val := range err.fields {
			if fields == "" {
				fields = fmt.Sprintf("%s:%v", key, val)
			} else {
				fields = fmt.Sprintf("%s\t%s:%v", fields, key, val)
			}
		}
	}
	return fields
}
