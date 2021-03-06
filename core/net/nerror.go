/********************************************************************************
* error.go
*
* Written by azraid@gmail.com
* Owned by azraid@gmail.com
********************************************************************************/

package net

import (
	"fmt"
	"os"
	"runtime"
)

type nerror struct {
	code int
	text string
}

type FErrorName func(code int) string

//RaiseNError(code, runtimeSkip = 1, data...)
func RaiseNError(ename FErrorName, args ...interface{}) NError {
	if len(args) == 0 {
		return nerror{code: NErrorSucess}
	}

	c := nerror{code: args[0].(int)}

	if len(args) > 1 {
		_, file, line, _ := runtime.Caller(args[1].(int))
		c.text = fmt.Sprintf("%s;%s;%s(%d);", ename(c.code), os.Args[0], file, line)
	}

	if l := len(args); l > 2 {
		for i := 2; i < l; i++ {
			c.text += fmt.Sprintf("%+v;", args[i])
		}
	}
	return c
}

func (e nerror) Error() string {
	return e.text
}

func (e nerror) Code() int {
	return e.code
}

func (e nerror) IsSuccess() bool {
	return e.code == NErrorSucess
}

func Sucess() NError {
	return nerror{code: NErrorSucess, text: "Sucess"}
}

const (
	NErrorSucess          = 0
	NErrorParsingError    = 1
	NErrorNotImplemented  = 2
	NErrorFederationError = 3
	NErrorAppStopping     = 4
	NErrorTooLargeSize    = 5
	NErrorUnknownMsgType  = 6
	NErrorInternal        = 7
	NErrorTimeout         = 8
	NErrorInvalidparams   = 9
	NErrorNoPermission    = 10
)

func CoErrorName(code int) string {
	switch code {
	case NErrorSucess:
		return "Sucess"
	case NErrorParsingError:
		return "NErrorParsingError"
	case NErrorNotImplemented:
		return "NErrorNotImplemented"
	case NErrorFederationError:
		return "NErrorFederationError"
	case NErrorAppStopping:
		return "NErrorAppStopping"
	case NErrorTooLargeSize:
		return "NErrorTooLargeSize"
	case NErrorUnknownMsgType:
		return "NErrorUnknownMsgType"
	case NErrorInternal:
		return "NErrorInternal"
	case NErrorTimeout:
		return "NErrorTimeout"
	case NErrorInvalidparams:
		return "NErrorInvalidparams"
	case NErrorNoPermission:
		return "NErrorNoPermission"
	}

	return "NErrorUnknown"
}

func CoRaiseNError(args ...interface{}) NError {
	return RaiseNError(CoErrorName, args)
}
