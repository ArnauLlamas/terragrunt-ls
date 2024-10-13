package lsp

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/hcl"
	lsp "go.lsp.dev/protocol"

	_ "github.com/ArnauLlamas/terragrunt-ls/internal/log"
)

func ParseAst(content []byte) *sitter.Tree {
	parser := sitter.NewParser()
	parser.SetLanguage(hcl.GetLanguage())
	tree, _ := parser.ParseCtx(context.Background(), nil, content)
	return tree
}

func (d *Document) ApplyChangesToAst(newContent string) {
	d.Ast = ParseAst([]byte(newContent))
}

func NodeAtPosition(tree *sitter.Tree, position lsp.Position) *sitter.Node {
	start := sitter.Point{Row: position.Line, Column: position.Character}
	return tree.RootNode().NamedDescendantForPointRange(start, start)
}

// Util function for debugging mainly
func LogNode(node *sitter.Node, content string) {
	log.Error(fmt.Sprintf("%s ; [%v, %v] - [%v, %v] ; %v",
		node.Type(),
		node.StartPoint().Row,
		node.StartPoint().Column,
		node.EndPoint().Row,
		node.EndPoint().Column,
		node.Content([]byte(content)),
	))
}

// func TreeCursor(n *sitter.Node) *sitter.TreeCursor {
// 	return sitter.NewTreeCursor(n)
// }
//
// func TreeCursorRoot(tree *sitter.Tree) *sitter.TreeCursor {
// 	return sitter.NewTreeCursor(tree.RootNode())
// }
