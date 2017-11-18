/********************************************************************************
* gridctx.go
* grid context에 대한 관리
*
* Written by azraid@gmail.com (2017-02-10)
* Owned by azraid@gmail.com
********************************************************************************/

package util

import "hash/fnv"

func Hash32(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
