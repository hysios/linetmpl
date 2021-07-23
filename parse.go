package linetmpl

import "text/template/parse"

const (
	LeftDelim  = "["
	RightDelim = "]"
)

func Parse(name, text string) (*Tree, error) {
	ast, err := parse.Parse(name, text, LeftDelim, RightDelim, Funcs)
	if err != nil {
		return nil, err
	}

	var tr = Traversal{Root: ast[name].Root}
	if tree, err := tr.Visitor(func(node parse.Node, child *Node) error {
		return nil
	}); err != nil {
		return nil, err
	} else {
		return tree, nil
	}
}
