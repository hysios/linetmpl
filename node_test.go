package linetmpl

import (
	"encoding/json"
	"testing"
	"text/template/parse"

	"github.com/tj/assert"
)

func testTree(text string) *parse.Tree {
	tree, _ := parse.Parse("test", text, "[", "]", Funcs)
	return tree["test"]
}

func TestTraversal_Visitor(t *testing.T) {
	t1 := testTree(`/[.日期]/[.车牌]/[.年]/[.月]/[.日]/[.违法代码]_[.序号].jpg`)
	tr := Traversal{Root: t1.Root}
	out, err := tr.Visitor(func(node parse.Node, child *Node) error {
		return nil
	})

	assert.NoError(t, err)
	printJ(t, out)
}

func printJ(t *testing.T, val interface{}) {
	b, err := json.MarshalIndent(val, "", "    ")
	if err != nil {
		t.Fatalf("error %s", err)
	}
	t.Logf("%s", string(b))
}
