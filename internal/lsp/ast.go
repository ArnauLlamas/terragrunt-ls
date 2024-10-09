package lsp

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/hcl"
	lsp "go.lsp.dev/protocol"
)

func ParseAst(content []byte) *sitter.Tree {
	parser := sitter.NewParser()
	parser.SetLanguage(hcl.GetLanguage())
	tree, _ := parser.ParseCtx(context.Background(), nil, content)
	return tree
}

func NodeAtPosition(tree *sitter.Tree, position lsp.Position) *sitter.Node {
	start := sitter.Point{Row: position.Line, Column: position.Character}
	return tree.RootNode().NamedDescendantForPointRange(start, start)
}

func TreeCursor(n *sitter.Node) *sitter.TreeCursor {
	return sitter.NewTreeCursor(n)
}

func TreeCursorRoot(tree *sitter.Tree) *sitter.TreeCursor {
	return sitter.NewTreeCursor(tree.RootNode())
}

// func FindLocals(tree *sitter.Tree) *sitter.Node {
// 	localsBlockPattern := `(
// 		(identifier) @constant
// 		(#eq? @constant "locals")
// 	)`
// 	q, _ := sitter.NewQuery([]byte(localsBlockPattern), hcl.GetLanguage())
// 	qc := sitter.NewQueryCursor()
// 	qc.Exec(q, tree.RootNode())
//
// 	for {
// 		m, ok := qc.NextMatch()
// 		if !ok {
// 			break
// 		}
//
// 		m = qc.FilterPredicates(m, )
// 	}
// }

func (d *Document) ApplyChangesToAst(newContent string) {
	d.Ast = ParseAst([]byte(newContent))
}
