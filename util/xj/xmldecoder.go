package xj

import (
	"encoding/xml"
	"io"
	"unicode"

	"golang.org/x/net/html/charset"
)

const (
	attrPrefix    = "-"
	contentPrefix = "#"
)

const (
	decodeStateStart = iota
	decodeStateData
	decodeStateEnd
)

// A xmlDecoder reads and decodes XML objects from an input stream.
type xmlDecoder struct {
	r   io.Reader
	err error
}

// NewxmlDecoder returns a new decoder that reads from r.
func newXmlDecoder(r io.Reader) *xmlDecoder {
	return &xmlDecoder{r: r}
}

// Decode reads the next JSON-encoded value from its
// input and stores it in the value pointed to by v.
func (dec *xmlDecoder) Decode(root *node) error {

	xmlDec := xml.NewDecoder(dec.r)

	// That will convert the charset if the provided XML is non-UTF-8
	xmlDec.CharsetReader = charset.NewReaderLabel

	// Create first element from the root node
	elem := &element{
		parent: nil,
		n:      root,
	}

	state := decodeStateEnd

	for {
		t, err := xmlDec.Token()
		//fmt.Printf("%v\r\n", t)

		if err != nil && err != io.EOF {
			return err
		}

		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:

			// Build new a new current element and link it to its parent
			elem = &element{
				parent: elem,
				n:      &node{},
				label:  se.Name.Local,
			}

			// Extract attributes as children
			// for _, a := range se.Attr {
			// 	elem.n.addChild(attrPrefix+a.Name.Local, &node{Data: a.Value})
			// }

			state = decodeStateStart

		case xml.CharData:
			if state == decodeStateStart {
				// Extract XML data (if any)
				elem.n.Data = trimNonGraphic(string(xml.CharData(se)))
				state = decodeStateData
			}
		case xml.EndElement:
			// And add it to its parent list
			if elem.parent != nil {
				elem.parent.addChild(elem.label, elem.n)
			}
			// Then change the current element to its parent
			elem = elem.parent
			state = decodeStateEnd
			// fmt.Printf("%v", elem)
		}
	}

	return nil
}

// trimNonGraphic returns a slice of the string s, with all leading and trailing
// non graphic characters and spaces removed.
//
// Graphic characters include letters, marks, numbers, punctuation, symbols,
// and spaces, from categories L, M, N, P, S, Zs.
// Spacing characters are set by category Z and property Pattern_White_Space.
func trimNonGraphic(s string) string {
	if s == "" {
		return s
	}

	var first *int
	var last int
	for i, r := range []rune(s) {
		if !unicode.IsGraphic(r) || unicode.IsSpace(r) {
			continue
		}

		if first == nil {
			f := i // copy i
			first = &f
			last = i
		} else {
			last = i
		}
	}

	// If first is nil, it means there are no graphic characters
	if first == nil {
		return ""
	}

	return string([]rune(s)[*first : last+1])
}
