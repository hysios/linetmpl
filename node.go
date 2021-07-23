package linetmpl

import (
	"encoding/json"
	"text/template/parse"
)

//go:generate stringer -type NodeType
type NodeType int

const (
	TEXT NodeType = iota
	ACTION
	PIPE
	COMMAND
	VARIABLE
	IDENTIFIER
	DOT
	NIL
	FIELD
	CHAIN
	BOOL
	NUMBER
	STRING
)

type Tree struct {
	Name  string  `json:"name,omitempty"`
	Nodes []*Node `json:"nodes,omitempty"`
}

type Node struct {
	NodeType NodeType `json:"nodeType"`
	Pos      int      `json:"pos"`
	Line     int      `json:"line,omitempty"`
	Text     string   `json:"text,omitempty"`
	Pipe     *Node    `json:"pipe,omitempty"`
	IsAssign bool     `json:"isAssign,omitempty"`
	Ident    []string `json:"ident,omitempty"`
	Nodes    []*Node  `json:"nodes,omitempty"`
	Decl     []*Node  `json:"decl,omitempty"`
	Cmds     []*Node  `json:"cmds,omitempty"`
	Args     []*Node  `json:"args,omitempty"`
	Node     *Node    `json:"node,omitempty"`
	Field    []string `json:"field,omitempty"`
	True     bool     `json:"true,omitempty"`
}

type VisitFunc func(node parse.Node, child *Node) error

type Traversal struct {
	Root *parse.ListNode
}

func (trversal *Traversal) Visitor(fn VisitFunc) (*Tree, error) {
	var (
		tree = &Tree{}
		root = trversal.Root
	)
	for _, child := range root.Nodes {
		if err := trversal.visitBytree(child, tree, fn); err != nil {
			return nil, err
		}

	}
	return tree, nil
}

func (trversal *Traversal) visitBytree(node parse.Node, tree *Tree, fn VisitFunc) error {
	if n, err := trversal.visitNode(node, fn); err != nil {
		return err
	} else {
		tree.Nodes = append(tree.Nodes, n)
	}
	return nil
}

func (trversal *Traversal) visitNode(node parse.Node, fn VisitFunc) (*Node, error) {
	var n *Node
	switch x := node.(type) {
	case *parse.TextNode:
		n = &Node{NodeType: TEXT, Pos: int(x.Pos)}
		if err := trversal.visitText(x, n, fn); err != nil {
			return nil, err
		}
	case *parse.ActionNode:
		n = &Node{NodeType: ACTION, Pos: int(x.Pos)}
		if err := trversal.visitAction(x, n, fn); err != nil {
			return nil, err
		}
	case *parse.CommentNode:
		// skip
	case *parse.PipeNode:
		n = &Node{NodeType: PIPE, Pos: int(x.Pos)}
		if err := trversal.visitPipe(x, n, fn); err != nil {
			return nil, err
		}
	case *parse.CommandNode:
		n = &Node{NodeType: COMMAND, Pos: int(x.Pos)}
		if err := trversal.visitCommand(x, n, fn); err != nil {
			return nil, err
		}
	case *parse.IdentifierNode:
		n = &Node{NodeType: IDENTIFIER, Pos: int(x.Pos)}
		if err := trversal.visitIdentifier(x, n, fn); err != nil {
			return nil, err
		}
	case *parse.VariableNode:
		n = &Node{NodeType: VARIABLE, Pos: int(x.Pos)}
		if err := trversal.visitVariable(x, n, fn); err != nil {
			return nil, err
		}
	case *parse.DotNode:
		n = &Node{NodeType: DOT, Pos: int(x.Pos)}
		if err := trversal.visitDot(x, n, fn); err != nil {
			return nil, err
		}
	case *parse.NilNode:
		n = &Node{NodeType: NIL, Pos: int(x.Pos)}
		if err := trversal.visitNil(x, n, fn); err != nil {
			return nil, err
		}
	case *parse.FieldNode:
		n = &Node{NodeType: FIELD, Pos: int(x.Pos)}
		if err := trversal.visitField(x, n, fn); err != nil {
			return nil, err
		}
	case *parse.ChainNode:
		n = &Node{NodeType: CHAIN, Pos: int(x.Pos)}
		if err := trversal.visitChain(x, n, fn); err != nil {
			return nil, err
		}
	case *parse.BoolNode:
		n = &Node{NodeType: BOOL, Pos: int(x.Pos)}
		if err := trversal.visitBool(x, n, fn); err != nil {
			return nil, err
		}
	case *parse.NumberNode:
		n = &Node{NodeType: NUMBER, Pos: int(x.Pos)}
		if err := trversal.visitNumber(x, n, fn); err != nil {
			return nil, err
		}
	case *parse.StringNode:
		n = &Node{NodeType: STRING, Pos: int(x.Pos)}
		if err := trversal.visitString(x, n, fn); err != nil {
			return nil, err
		}
	}
	return n, nil
}

func (trversal *Traversal) visitText(node *parse.TextNode, n *Node, fn VisitFunc) error {
	n.Text = string(node.Text)
	return fn(node, n)
}

func (trversal *Traversal) visitComment(node *parse.CommentNode, n *Node, fn VisitFunc) error {
	return nil
	// panic("nonimplement")
}

// PipeNode
func (trversal *Traversal) visitPipe(node *parse.PipeNode, n *Node, fn VisitFunc) error {
	if node.IsAssign {
		n.IsAssign = true
	}

	if err := fn(node, n); err != nil {
		return err
	}

	for _, decl := range node.Decl {
		var nn = &Node{NodeType: VARIABLE, Pos: int(decl.Pos)}
		if err := trversal.visitVariable(decl, nn, fn); err != nil {
			return err
		}
		n.Decl = append(n.Decl, nn)
	}

	for _, cmd := range node.Cmds {
		var nn = &Node{NodeType: COMMAND, Pos: int(cmd.Pos)}
		if err := trversal.visitCommand(cmd, nn, fn); err != nil {
			return err
		}
		n.Cmds = append(n.Cmds, nn)
	}
	return nil
}

func (trversal *Traversal) visitAction(node *parse.ActionNode, n *Node, fn VisitFunc) error {
	if err := fn(node, n); err != nil {
		return err
	}

	var nn = &Node{NodeType: PIPE, Pos: int(node.Pos)}

	if err := trversal.visitPipe(node.Pipe, nn, fn); err != nil {
		return err
	}
	n.Pipe = nn
	return nil
}

// CommandNode
func (trversal *Traversal) visitCommand(node *parse.CommandNode, n *Node, fn VisitFunc) error {
	if err := fn(node, n); err != nil {
		return err
	}

	for _, arg := range node.Args {
		// nn := &Node{Pos: int(arg.Position())}
		if nn, err := trversal.visitNode(arg, fn); err != nil {
			return err
		} else {
			n.Args = append(n.Args, nn)
		}
	}
	return nil
}

// IdentifierNode
func (trversal *Traversal) visitIdentifier(node *parse.IdentifierNode, n *Node, fn VisitFunc) error {
	n.Ident = append(n.Ident, node.Ident)
	if err := fn(node, n); err != nil {
		return err
	}
	return nil
}

// VariableNode
func (trversal *Traversal) visitVariable(node *parse.VariableNode, n *Node, fn VisitFunc) error {
	n.Ident = node.Ident
	if err := fn(node, n); err != nil {

		return err
	}
	return nil
}

// DotNode
func (trversal *Traversal) visitDot(node *parse.DotNode, n *Node, fn VisitFunc) error {
	if err := fn(node, n); err != nil {
		return err
	}
	return nil
}

// NilNode
func (trversal *Traversal) visitNil(node *parse.NilNode, n *Node, fn VisitFunc) error {
	if err := fn(node, n); err != nil {
		return err
	}
	return nil
}

// FieldNode
func (trversal *Traversal) visitField(node *parse.FieldNode, n *Node, fn VisitFunc) error {
	n.Ident = node.Ident
	if err := fn(node, n); err != nil {

		return err
	}
	return nil
}

// ChainNode
func (trversal *Traversal) visitChain(node *parse.ChainNode, n *Node, fn VisitFunc) error {
	if err := fn(node, n); err != nil {
		return err
	}
	if node.Node != nil {
		if nn, err := trversal.visitNode(node.Node, fn); err != nil {
			return err
		} else {
			n.Node = nn
		}
	}

	n.Field = node.Field
	return nil
}

// BoolNode
func (trversal *Traversal) visitBool(node *parse.BoolNode, n *Node, fn VisitFunc) error {
	n.True = node.True
	if err := fn(node, n); err != nil {
		return err
	}
	return nil
}

// NumberNode
func (trversal *Traversal) visitNumber(node *parse.NumberNode, n *Node, fn VisitFunc) error {
	n.Text = node.Text
	if err := fn(node, n); err != nil {
		return err
	}
	return nil
}

// StringNode}
func (trversal *Traversal) visitString(node *parse.StringNode, n *Node, fn VisitFunc) error {
	n.Text = node.Text
	if err := fn(node, n); err != nil {
		return err
	}
	return nil
}

// BranchNode
func (trversal *Traversal) visitBranch(node *parse.BranchNode, n *Node, fn VisitFunc) error {
	panic("nonimplement")
}

// IfNode
func (trversal *Traversal) visitIf(node *parse.IfNode, n *Node, fn VisitFunc) error {
	panic("nonimplement")
}

// RangeNode
func (trversal *Traversal) visitRange(node *parse.RangeNode, n *Node, fn VisitFunc) error {
	panic("nonimplement")
}

// TemplateNode
func (trversal *Traversal) visitTemplate(node *parse.TemplateNode, n *Node, fn VisitFunc) error {
	panic("nonimplement")
}

func (n NodeType) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.String())
}
