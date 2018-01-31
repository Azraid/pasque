/********************************************************************************
* error.go
*
* Written by azraid@gmail.com (2016-07-26)
* Owned by azraid@gmail.com
********************************************************************************/

package core

import (
	"fmt"
	"runtime"
)

const (
	NetErrorSucess          = 0
	NetErrorParsingError    = 1
	NetErrorNotImplemented  = 2
	NetErrorFederationError = 3
	NetErrorAppStopping     = 4
	NetErrorTooLargeSize    = 5
	NetErrorUnknownMsgType  = 6
	NetErrorInternal        = 7
	NetErrorTimeout         = 8
	NetErrorInvalidparams   = 9
	NetErrorNoPermission    = 10
)

type NetError struct {
	Code uint32 `json:",string"`
	Text string
}

func (nerr NetError) Error() string {
	return fmt.Sprintf("code:%d, text:%s, issue:%s", nerr.Code, nerr.Text)
}

func Error(code uint32, text string) NetError {
	_, file, line, _ := runtime.Caller(1)
	return NetError{Code: code, Text: text + fmt.Sprintf("; %s(%d)", file, line)}
}

func Sucess() NetError {
	return NetError{Code: NetErrorSucess}
}
