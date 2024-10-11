package handler

import (
	"context"
	"encoding/json"

	"go.lsp.dev/jsonrpc2"
	lsp "go.lsp.dev/protocol"

	"github.com/ArnauLlamas/terragrunt-ls/internal/docs"
	_ "github.com/ArnauLlamas/terragrunt-ls/internal/log"
)

var completionItems []lsp.CompletionItem

func init() {
	// TODO: Split into different slices so we can return only the ones
	// that make sense in context
	buildFunctions(&completionItems)
	buildBlocks(&completionItems)
	buildAttributes(&completionItems)
	buildTopLevelAttributes(&completionItems)
}

func (h *handler) completion(
	ctx context.Context,
	reply jsonrpc2.Replier,
	req jsonrpc2.Request,
) error {
	var err error
	var params lsp.CompletionParams
	// var completionItems []lsp.CompletionItem

	err = json.Unmarshal(req.Params(), &params)
	if err != nil {
		return err
	}

	// TODO: Logic based on tree location
	// doc := h.documents.GetDocument(params.TextDocument.URI)
	// if doc.Ast.RootNode().HasChanges() {
	// 	doc.ApplyChangesToAst(doc.Content)
	// }

	// currentNode := lsplocal.NodeAtPosition(doc.Ast, params.Position)
	// getLocals(&completionItems, doc.Ast, []byte(doc.Content))

	return reply(ctx, completionItems, nil)
}

func createMarkupContent(doc string) lsp.MarkupContent {
	markupDoc := lsp.MarkupContent{
		Kind:  lsp.MarkupKind("markdown"),
		Value: doc,
	}
	return markupDoc
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
