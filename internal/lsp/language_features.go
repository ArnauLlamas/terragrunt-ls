package lsp

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/hcl"
)

var HCL *sitter.Language = hcl.GetLanguage()

func IsNodeAtTopLevel(tree *sitter.Tree, node *sitter.Node) bool {
	return node == tree.RootNode() || node.Parent() == tree.RootNode().Child(0)
}

func GetLocals(tree *sitter.Tree, content string) []string {
	pattern := `(
		(identifier) @constant
		(#match? @constant "^locals$")
	)`

	return getBlocksAttributes(tree, []byte(pattern), []byte(content))
}

func GetIncludes(tree *sitter.Tree, content string) []string {
	pattern := `(
		(identifier) @constant
		(#match? @constant "^include$")
	)`

	return getBlocksLabels(tree, []byte(pattern), []byte(content))
}

func GetDependencys(tree *sitter.Tree, content string) []string {
	pattern := `(
		(identifier) @constant
		(#match? @constant "^dependency$")
	)`

	return getBlocksLabels(tree, []byte(pattern), []byte(content))
}

func getBlocksAttributes(tree *sitter.Tree, blockPattern []byte, content []byte) []string {
	q, _ := sitter.NewQuery(blockPattern, HCL)
	qc := sitter.NewQueryCursor()
	qc.Exec(q, tree.RootNode())

	var attrs []string
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}

		m = qc.FilterPredicates(m, content)
		for _, c := range m.Captures {
			attr := c.Node.NextSibling().NextSibling()
			for i := 0; i < int(attr.ChildCount()); i++ {
				attrs = append(attrs, attr.Child(i).Child(0).Content(content))
			}
		}
	}

	return attrs
}

func getBlocksLabels(tree *sitter.Tree, blockPattern []byte, content []byte) []string {
	q, _ := sitter.NewQuery(blockPattern, HCL)
	qc := sitter.NewQueryCursor()
	qc.Exec(q, tree.RootNode())

	var labels []string
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}

		m = qc.FilterPredicates(m, content)
		for _, c := range m.Captures {
			label := c.Node.NextSibling().Child(1)
			labels = append(labels, label.Content(content))
		}
	}

	return labels
}
