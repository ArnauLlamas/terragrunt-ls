package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	sitter "github.com/smacker/go-tree-sitter"
	"go.lsp.dev/jsonrpc2"
	lsp "go.lsp.dev/protocol"

	"github.com/ArnauLlamas/terragrunt-ls/internal/docs"
	_ "github.com/ArnauLlamas/terragrunt-ls/internal/log"
	lsplocal "github.com/ArnauLlamas/terragrunt-ls/internal/lsp"
)

var (
	blocks        []lsp.CompletionItem
	topLevelAttrs []lsp.CompletionItem
	attrs         []lsp.CompletionItem
	functions     []lsp.CompletionItem
)

var completionItems []lsp.CompletionItem

func init() {
	buildBlocks(&blocks)
	buildTopLevelAttributes(&topLevelAttrs)
	buildAttributes(&attrs)
	buildFunctions(&functions)
}

func (h *handler) completion(
	ctx context.Context,
	reply jsonrpc2.Replier,
	req jsonrpc2.Request,
) error {
	var err error
	var params lsp.CompletionParams

	err = json.Unmarshal(req.Params(), &params)
	if err != nil {
		return err
	}

	log.Debug("Looking for completions...")

	doc := h.documents.GetDocument(params.TextDocument.URI)

	// TODO: Return completions based on AST node position/type
	// Does not work right now but serves as inspiration
	// currentNode := lsplocal.NodeAtPosition(doc.Ast, params.Position)
	// if lsplocal.IsNodeAtTopLevel(doc.Ast, currentNode) {
	// 	completionItems = append(completionItems, blocks...)
	// 	completionItems = append(completionItems, topLevelAttrs...)
	//
	// 	return reply(ctx, completionItems, nil)
	//
	// } else {
	// 	completionItems = append(completionItems, topLevelAttrs...)
	// 	completionItems = append(completionItems, functions...)
	//
	// 	buildLocals(&completionItems, doc.Ast, doc.Content)
	// 	buildIncludes(&completionItems, doc.Ast, doc.Content)
	// 	buildDependencys(&completionItems, doc.Ast, doc.Content)
	//
	// 	return reply(ctx, completionItems, nil)
	// }

	// TODO: Return attributes completions of current block
	// Does not work right now but serves as inspiration
	// lsplocal.LogNode(currentNode, doc.Content)
	// lsplocal.LogNode(currentNode.Parent(), doc.Content)
	//
	// if currentNode.Type() == "ERROR" &&
	// 	strings.HasPrefix(currentNode.Content([]byte(doc.Content)), "dependency") {
	// }
	//
	// return reply(ctx, completionItems, nil)

	// WARN: Just return all completions for now
	completionItems = append(completionItems, blocks...)
	completionItems = append(completionItems, topLevelAttrs...)
	completionItems = append(completionItems, topLevelAttrs...)
	completionItems = append(completionItems, functions...)

	// following function should also be retriggered on save doc/change doc
	h.updateVariables(&completionItems, doc.Ast, doc.Content, params.Position)
	h.removeDuplicatesCompletionItems(&completionItems)

	return reply(ctx, completionItems, nil)
}

func (h *handler) removeDuplicatesCompletionItems(items *[]lsp.CompletionItem) {
	keys := make(map[string]bool)
	list := []lsp.CompletionItem{}

	for _, item := range *items {
		if _, value := keys[item.Label]; !value {
			keys[item.Label] = true
			list = append(list, item)
		}
	}

	*items = list
}

func (h *handler) updateVariables(
	items *[]lsp.CompletionItem,
	tree *sitter.Tree,
	content string,
	position lsp.Position,
) {
	buildLocals(items, tree, content, position)
	buildIncludes(items, tree, content)
	buildDependencys(items, tree, content)
}

func createMarkupContent(doc string) lsp.MarkupContent {
	markupDoc := lsp.MarkupContent{
		Kind:  lsp.MarkupKind("markdown"),
		Value: doc,
	}
	return markupDoc
}

func deleteCompletionItemsThatHasPrefix(items *[]lsp.CompletionItem, prefix string) {
	list := []lsp.CompletionItem{}
	for _, item := range *items {
		if !strings.HasPrefix(item.Label, prefix) {
			list = append(list, item)
		}
	}

	*items = list
}

func buildLocals(
	items *[]lsp.CompletionItem,
	tree *sitter.Tree,
	content string,
	position lsp.Position,
) {
	// Quick exit if we are in locals block
	currentNode := lsplocal.NodeAtPosition(tree, position)
	if currentNode.Content([]byte(content)) == "locals" {
		return
	}

	deleteCompletionItemsThatHasPrefix(items, "local.")

	locals := lsplocal.GetLocals(tree, content)
	for _, local := range locals {
		item := lsp.CompletionItem{
			Kind:       lsp.CompletionItemKindVariable,
			Label:      fmt.Sprintf("local.%s", local),
			Detail:     fmt.Sprintf("local.%s", local),
			InsertText: fmt.Sprintf("local.%s", local),
		}

		*items = append(*items, item)
	}
}

func buildIncludes(items *[]lsp.CompletionItem, tree *sitter.Tree, content string) {
	deleteCompletionItemsThatHasPrefix(items, "include.")

	includes := lsplocal.GetIncludes(tree, content)

	for _, include := range includes {
		item := lsp.CompletionItem{
			Kind:       lsp.CompletionItemKindVariable,
			Label:      fmt.Sprintf("include.%s", include),
			Detail:     fmt.Sprintf("include.%s", include),
			InsertText: fmt.Sprintf("include.%s", include),
		}

		*items = append(*items, item)
	}
}

func buildDependencys(items *[]lsp.CompletionItem, tree *sitter.Tree, content string) {
	deleteCompletionItemsThatHasPrefix(items, "dependency.")

	dependencys := lsplocal.GetDependencys(tree, content)

	for _, dependency := range dependencys {
		item := lsp.CompletionItem{
			Kind:       lsp.CompletionItemKindVariable,
			Label:      fmt.Sprintf("dependency.%s", dependency),
			Detail:     fmt.Sprintf("dependency.%s", dependency),
			InsertText: fmt.Sprintf("dependency.%s", dependency),
		}

		*items = append(*items, item)
	}
}

func buildFunctions(items *[]lsp.CompletionItem) {
	functions := docs.GetFunctions()

	for _, function := range functions {
		item := lsp.CompletionItem{
			Kind:             lsp.CompletionItemKindFunction,
			Label:            function.Item,
			Detail:           function.InsertText,
			InsertText:       function.InsertText,
			InsertTextFormat: lsp.InsertTextFormatSnippet,
			Documentation:    createMarkupContent(function.Content),
		}

		*items = append(*items, item)
	}
}

func buildBlocks(items *[]lsp.CompletionItem) {
	blocks := docs.GetBlocks()

	for _, block := range blocks {
		item := lsp.CompletionItem{
			Kind:             lsp.CompletionItemKindSnippet,
			Label:            block.Item,
			Detail:           block.Item,
			InsertText:       block.InsertText,
			InsertTextFormat: lsp.InsertTextFormatSnippet,
			Documentation:    createMarkupContent(block.Content),
		}

		*items = append(*items, item)
	}
}

func buildTopLevelAttributes(items *[]lsp.CompletionItem) {
	attributes := docs.GetTopLevelAttributes()

	for _, attribute := range attributes {
		item := lsp.CompletionItem{
			Kind:             lsp.CompletionItemKindProperty,
			Label:            attribute.Item,
			Detail:           attribute.Item,
			InsertText:       attribute.InsertText,
			InsertTextFormat: lsp.InsertTextFormatSnippet,
			Documentation:    createMarkupContent(attribute.Content),
		}

		*items = append(*items, item)
	}
}

func buildAttributes(items *[]lsp.CompletionItem) {
	attributes := docs.GetAttributes()

	for _, attribute := range attributes {
		item := lsp.CompletionItem{
			Kind:             lsp.CompletionItemKindProperty,
			Label:            attribute.Item,
			Detail:           attribute.Item,
			InsertText:       attribute.InsertText,
			InsertTextFormat: lsp.InsertTextFormatSnippet,
			Documentation:    createMarkupContent(attribute.Content),
		}

		*items = append(*items, item)
	}
}
