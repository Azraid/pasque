package util

import "math/rand"

type RandSet interface {
	Add(value interface{})
	AnyOne() interface{}
	Range(f func(v interface{}) bool)
	Length() int
}

type randSet struct {
	values []interface{}
}

func NewRandSet() RandSet {
	return &randSet{}
}

func (s *randSet) Length() int {
	return len(s.values)
}
func (s *randSet) Add(value interface{}) {
	s.values = append(s.values, value)
}

func (s *randSet) AnyOne() interface{} {
	return s.values[rand.Intn(len(s.values))]
}

func (s *randSet) Range(f func(v interface{}) bool) {
	for _, v := range s.values {
		if !f(v) {
			return
		}
	}
}
