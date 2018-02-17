/********************************************************************************
* util.go
*
* Written by azraid@gmail.com (2016-07-26)
* Owned by azraid@gmail.com
********************************************************************************/

package core

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const TIME_LAYOUT = "2006-01-02 15:04:05.999"

type TimeNZone time.Time
type TUserID string

func ToUserID(v string) TUserID {
	return TUserID(v)
}

func (t TimeNZone) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TIME_LAYOUT)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TIME_LAYOUT)
	b = append(b, '"')
	return b, nil
}

func (t *TimeNZone) UnmarshalJSON(b []byte) error {
	if tv, err := time.Parse(`"`+TIME_LAYOUT+`"`, string(b)); err != nil {
		return err
	} else {
		*t = TimeNZone(tv)
		return nil
	}
}

func (t TimeNZone) Sub(d time.Time) time.Duration {
	return time.Time(t).Sub(d)
}
func (t *TimeNZone) Parse(value string) error {
	if tv, err := time.Parse(TIME_LAYOUT, value); nil != err {
		return err
	} else {
		*t = TimeNZone(tv)
		return nil
	}
}
func (t TimeNZone) String() string {
	return time.Time(t).Format(TIME_LAYOUT)
}
func (t TimeNZone) IsPassed() bool {
	now := Now()
	dt := time.Time(t)
	if now.Year() > dt.Year() {
		return true
	} else if now.Year() < dt.Year() {
		return false
	}
	if now.YearDay() > dt.YearDay() {
		return true
	} else if now.YearDay() < dt.YearDay() {
		return false
	}
	if now.Hour() > dt.Hour() {
		return true
	} else if now.Hour() < dt.Hour() {
		return false
	}

	if now.Minute() > dt.Hour() {
		return true
	} else if now.Minute() < dt.Hour() {
		return false
	}

	if now.Second() > dt.Second() {
		return true
	}
	return false

}

type Guid [16]byte

const GuidSIZE = 38

var GuidZERO Guid

var re = regexp.MustCompile("^([A-z0-9]{8})-([A-z0-9]{4})-([A-z0-9]{4})-([A-z0-9]{4})-([A-z0-9]{12})$")

func ParseGuid(s string) (g Guid, err error) {
	md := re.FindStringSubmatch(s)
	if md == nil {
		return g, errors.New("Invalid Guid string")
	}
	b, err := hex.DecodeString(md[1] + md[2] + md[3] + md[4] + md[5])
	if err != nil {
		return g, errors.New("Invalid Guid string")
	}

	copy(g[:], b)
	return g, nil
}

func GenerateGuid() Guid {
	var g Guid

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, time.Now().UTC().UnixNano())
	copy(g[:8], buf.Bytes())

	for i := 0; i < 8; i++ {
		g[i+8] = byte(rand.Intn(16))
	}
	g[6] = (g[6] & 0xF) | (4 << 4)
	g[8] = (g[8] | 0x40) & 0x7F
	return g
}

func (g Guid) String() string {
	return fmt.Sprintf("%X-%X-%X-%X-%X", g[:4], g[4:6], g[6:8], g[8:10], g[10:])
}

func (g *Guid) MarshalJSON() ([]byte, error) {

	if g == nil {
		return nil, nil
	}
	b := make([]byte, 0, 36+2)
	b = append(b, '"')
	b = append(b, []byte(g.String())...)
	b = append(b, '"')

	return b, nil
}
func (g *Guid) UnmarshalJSON(b []byte) error {
	if GuidSIZE != len(b) {
		*g = GuidZERO
		return nil
	}

	if n, err := ParseGuid(string(b[1:37])); err != nil {
		return err
	} else {
		*g = n
	}
	return nil
}

type TRange struct {
	Min uint64
	Max uint64
}

type TStructTag struct {
	Required      bool
	LengthChecked bool
	RangeChecked  bool
	Length        int
	Range         TRange
}

func (p *TStructTag) ValidString(v reflect.Value, t reflect.StructField) error {
	str := v.String()
	strLen := len(str)
	if p.Required && strLen == 0 {
		return fmt.Errorf("%s required", t.Name)
	}
	if p.LengthChecked && strLen > p.Length {
		return fmt.Errorf("%s is greater than the length(%d)", t.Name, p.Length)
	}
	return nil
}
func (p *TStructTag) ValidPtr(v reflect.Value, t reflect.StructField) error {
	switch v.Type() {
	case reflect.TypeOf((*Guid)(nil)):
		if p.Required {
			if v.IsNil() {
				return fmt.Errorf("%s required", t.Name)
			}
			guid := v.Interface().(*Guid)
			if *guid == GuidZERO {
				return fmt.Errorf("%s required", t.Name)
			}
		}

	default:
		panic(fmt.Errorf("ValidPtr... undefined reflect.Type(%s)\n to do define type.go", v.Type()))
	}
	return nil
}

func (p *TStructTag) ValidInt(v reflect.Value, t reflect.StructField) error {
	val := v.Uint()
	if p.Required && val == 0 {
		return fmt.Errorf("%s required", t.Name)
	}
	if !p.Required && val == 0 {
		p.RangeChecked = false
	}
	if p.RangeChecked {
		if p.Range.Min > 0 && val < p.Range.Min {
			return fmt.Errorf("%s out of range min(%d)", t.Name, p.Range.Min)
		}
		if p.Range.Max > 0 && val > p.Range.Max {
			return fmt.Errorf("%s out of range max(%d)", t.Name, p.Range.Max)
		}
	}
	return nil
}
func (p *TStructTag) ValidArray(v reflect.Value, t reflect.StructField) error {
	if p.Required && (v.IsNil() || 0 == v.Len()) {
		return fmt.Errorf("%s required", t.Name)
	}
	if p.LengthChecked && v.Len() > p.Length {
		return fmt.Errorf("%s is greater than the length(%d)", t.Name, p.Length)
	}
	return nil
}
func setStructTag(tag reflect.StructTag) *TStructTag {
	p := &TStructTag{}
	p.Required = 0 == strings.Compare(tag.Get("required"), "true")
	_tLen := tag.Get("length")
	_tRange := tag.Get("range")
	if len(_tLen) > 0 {
		p.LengthChecked = true
		var err error
		if p.Length, err = strconv.Atoi(_tLen); nil != err {
			p.LengthChecked = false
		}
	}
	if len(_tRange) > 0 {
		p.RangeChecked = true
		if err := json.Unmarshal([]byte(_tRange), &p.Range); nil != err {
			p.RangeChecked = false
		}
	}
	return p
}

func CheckParam(v interface{}) *NetError {
	r := reflect.ValueOf(v)
	t := reflect.TypeOf(v)
	if r.Kind() == reflect.Ptr {
		panic(fmt.Errorf("ptr value is not allowed.\n change your code"))
	}
	var nErr *NetError
	var err error

	for i := 0; i < r.NumField(); i++ {
		tag := setStructTag(t.Field(i).Tag)
		switch r.Field(i).Kind() {
		case reflect.String:
			err = tag.ValidString(r.Field(i), t.Field(i))
		case reflect.Uint16:
			err = tag.ValidInt(r.Field(i), t.Field(i))
		case reflect.Uint32:
			err = tag.ValidInt(r.Field(i), t.Field(i))
		case reflect.Uint64:
			err = tag.ValidInt(r.Field(i), t.Field(i))
		case reflect.Ptr:
			err = tag.ValidPtr(r.Field(i), t.Field(i))
		case reflect.Slice:
			err = tag.ValidArray(r.Field(i), t.Field(i))
		case reflect.Array:
			err = tag.ValidArray(r.Field(i), t.Field(i))
		default:
			panic(fmt.Errorf("undefined Type For CheckParam\n to do define type.go"))
		}
		if nil != err {
			nErr = &NetError{Code: NetErrorInvalidparams, Text: err.Error()}
			break
		}
	}
	return nErr
}
