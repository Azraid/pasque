/*
*
* 이 코드는 우선 github.com/basgys/goxml2json를 사용한다.
* 향후에 수정하기로 하고..
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

// // addChild appends a node to the list of children
// func (n *node) addChild(s string, c *node) {
// 	// Lazy lazy
// 	if n.Children == nil {
// 		n.Children = map[string]nodes{}
// 	}

// 	n.Children[s] = append(n.Children[s], c)
// 	if len(n.Children[s]) > 1 {
// 	}
// }

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
