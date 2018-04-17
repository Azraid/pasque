package core

import (
	"fmt"
	"runtime"
)



func IssueErrorf(f string, v ...interface{}) error {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Errorf("%s;%s:%d", fmt.Sprintf(f, v), file, line)
}
