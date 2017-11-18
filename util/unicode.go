/********************************************************************************
* util.go
*
* Written by azraid@gmail.com (2016-07-26)
* Owned by azraid@gmail.com
********************************************************************************/

package util

import (
	"unicode/utf16"
	"unicode/utf8"
)

func WChar2String(w []uint16) string {
	for i, v := range w {
		if v == 0 {
			w = w[0:i]
			break
		}
	}

	return string(utf16.Decode(w))
}

func String2WChar(s string) []uint16 {
	var r8 []rune

	for i := 0; len(s) > 0; i++ {
		u, size := utf8.DecodeRuneInString(s)
		r8 = append(r8, u)
		s = s[size:]
	}

	r := utf16.Encode(r8)
	r = append(r, 0)

	return r
}
