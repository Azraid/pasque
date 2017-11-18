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
	NetErrorTimeout         = 42
	NetErrorInvalidparams   = 43
)

type NetError struct {
	Code  uint32 `json:",string"`
	Text  string
	Issue string
}

func (nerr NetError) Error() string {
	return fmt.Sprintf("code:%d, text:%s, issue:%s", nerr.Code, nerr.Text, nerr.Issue)
}

func Error(code uint32, text string) NetError {
	_, file, line, _ := runtime.Caller(1)

	return NetError{Code: code, Text: text, Issue: fmt.Sprintf("%s(%d)", file, line)}
}

func Sucess() NetError {
	return NetError{Code: NetErrorSucess}
}
