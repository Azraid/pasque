/*
*
 */

package xj

import "bytes"

// node is a data element on a tree
type node struct {
	Children map[string]nodes
	Data     string
}

type element struct {
	parent *element
	n      *node
	label  string
}

// nodes is a list of nodes
type nodes []*node

func (e *element) addChild(s string, c *node) {
	// Lazy lazy
	if e.n.Children == nil {
		e.n.Children = map[string]nodes{}
	}

	e.n.Children[s] = append(e.n.Children[s], c)
}

// isComplex returns whether it is a complex type (has children)
func (n *node) isComplex() bool {
	return len(n.Children) > 0
}

func Xml2Json(x []byte) ([]byte, error) {
	if len(x) == 0 {
		return []byte("{}"), nil
	}

	root := &node{}
	err := newXmlDecoder(bytes.NewReader(x)).Decode(root)
	if err != nil {
		return nil, err
	}

	// Then encode it in JSON
	buf := new(bytes.Buffer)
	err = newJsonEncoder(buf).Encode(root)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
