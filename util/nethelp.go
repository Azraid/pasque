/********************************************************************************
* util.go
*
* Written by azraid@gmail.com (2016-07-26)
* Owned by azraid@gmail.com
********************************************************************************/

package util

import (
	"net"
	"strings"
)

func LocalIP4(conn net.Conn) net.IP {
	if conn != nil {
		ip4 := strings.Split(conn.LocalAddr().String(), ":")[0]
		return net.ParseIP(ip4).To4()
	}
	return nil
}

func StrCmpI(a1 string, a2 string) bool {
	return strings.EqualFold(a1, a2)
}
