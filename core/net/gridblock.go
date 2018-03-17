/********************************************************************************
* gridblock.go
*
* Written by azraid@gmail.com
* Owned by azraid@gmail.com
********************************************************************************/

package net

import (
	"math/rand"
	"sync"

	. "github.com/Azraid/pasque/core"
	"github.com/Azraid/pasque/util"
)

type FederatedApi struct {
	lock *sync.RWMutex
	apis []string
	key  string
}

func NewFederatedApi() *FederatedApi {
	return &FederatedApi{lock: new(sync.RWMutex)}
}

//AssiginKey는 딱 하나만 할 당된다.
func (fa *FederatedApi) AssignKey(key string) bool {
	if len(fa.key) == 0 {
		fa.key = key
		return true
	}

	return (fa.key == key)
}

func (fa *FederatedApi) Register(api string) error {
	fa.lock.Lock()
	defer fa.lock.Unlock()

	for _, v := range fa.apis {
		if util.StrCmpI(v, api) {
			return nil
		}
	}

	fa.apis = append(fa.apis, api)
	return nil
}

func (fa FederatedApi) Find(api string) bool {
	for _, v := range fa.apis {
		if util.StrCmpI(v, api) {
			return true
		}
	}

	return false
}

func (fa FederatedApi) Compare(key string) bool {
	return util.StrCmpI(key, fa.key)
}

//GridBlock 의 buckets[] "eid1", "eid2", "eid3", "eid4"
//GridBlock의  buckets은 처음 서버가 구동될때, config에 의 해서 buckets을 할당한다.
//config에서 할당을 끝내면, fixup으로 그 크기를 고정해야 한다.
//이는 key distribution 알고리즘이 고정 hash키를 사용하기 때문이다.
type GridBlock struct {
	buckets []string
	max     uint32
	lock    *sync.RWMutex
}

func NewGridBlock() *GridBlock {
	return &GridBlock{lock: new(sync.RWMutex)}
}

func (gb *GridBlock) Register(eid string) error {
	gb.lock.Lock()
	defer gb.lock.Unlock()

	for _, v := range gb.buckets {
		if util.StrCmpI(v, eid) {
			return IssueErrorf("%s already exists", eid)
		}
	}

	gb.buckets = append(gb.buckets, eid)

	return nil
}

func (gb *GridBlock) Fixup() {
	gb.max = uint32(len(gb.buckets))
}

func (gb *GridBlock) Distribute(key string) string {
	if len(key) > 0 {
		h := util.Hash32(key)
		return gb.buckets[h%gb.max]
	}

	return gb.buckets[rand.Intn(int(gb.max))]
}
