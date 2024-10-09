package handler

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"go.lsp.dev/jsonrpc2"
	lsp "go.lsp.dev/protocol"

	_ "github.com/ArnauLlamas/terragrunt-ls/internal/log"
	// lsplocal "github.com/ArnauLlamas/terragrunt-ls/internal/lsp"
)

const BASE_DOCS_PATH string = "terragrunt"

//go:embed terragrunt/**/*.md
var documentation embed.FS

var completionItems []lsp.CompletionItem

func init() {
	// TODO: Split into different slices so we can return only the ones
	// that make sense in context
	getFunctions(&completionItems)
	getBlocks(&completionItems)
	getAttributes(&completionItems)
}

func (h *handler) completion(
	ctx context.Context,
	reply jsonrpc2.Replier,
	req jsonrpc2.Request,
) error {
	var err error
	var params lsp.CompletionParams
	// var completionItemsTwo []lsp.CompletionItem

	err = json.Unmarshal(req.Params(), &params)
	if err != nil {
		return err
	}

	// TODO: Logic based on tree location
	// doc := h.documents.GetDocument(params.TextDocument.URI)
	// currentNode := lsplocal.NodeAtPosition(doc.Ast, params.Position)

	return reply(ctx, completionItems, nil)
}

func createMarkupContent(doc string) lsp.MarkupContent {
	markupDoc := lsp.MarkupContent{
		Kind:  lsp.MarkupKind("markdown"),
		Value: doc,
	}
	return markupDoc
}

func getDocumentationContentsBasedOnDirName(dirName string) (docs []string, err error) {
	docsPath := fmt.Sprintf("%s/%s", BASE_DOCS_PATH, dirName)
	dirEntries, err := documentation.ReadDir(docsPath)
	if err != nil {
		errorMsg := fmt.Sprintf("Cannot read documentation on %s", docsPath)
		log.Error(errorMsg)
		return nil, err
	}

	for _, entry := range dirEntries {
		docFile := entry.Name()
		docContent := readEmbeddedFile(fmt.Sprintf("%s/%s", docsPath, docFile))

		docs = append(docs, docContent)
	}

	return docs, nil
}

func readEmbeddedFile(fileName string) (fileContent string) {
	filePath := fmt.Sprintf(fileName)
	fileBytes, _ := documentation.ReadFile(filePath)

	return string(fileBytes)
}

func getFunctions(items *[]lsp.CompletionItem) {
	docs, err := getDocumentationContentsBasedOnDirName("functions")
	if err != nil {
		log.Panic("Failed to read functions documentation")
		panic(1)
	}

	for _, doc := range docs {
		docLines := strings.Split(doc, "\n")

		functionSignature := docLines[0]
		functionName := strings.Split(functionSignature, "(")[0]
		content := strings.Join(docLines[1:], "\n")

		item := lsp.CompletionItem{
			Kind:             lsp.CompletionItemKindFunction,
			Label:            functionName,
			Detail:           functionSignature,
			InsertText:       functionSignature,
			InsertTextFormat: lsp.InsertTextFormatSnippet,
			Documentation:    createMarkupContent(content),
		}

		*items = append(*items, item)
	}
}

func getBlocks(items *[]lsp.CompletionItem) {
	docs, err := getDocumentationContentsBasedOnDirName("blocks")
	if err != nil {
		log.Panic("Failed to read blocks documentation")
		panic(1)
	}

	for _, doc := range docs {
		docLines := strings.Split(doc, "\n")

		blockName := docLines[0]
		content := strings.Join(docLines[1:], "\n")

		// Some blocks are named, these ones will have a $name string
		// before the opened block in the InsertText field
		var insertText string
		switch blockName {
		case "dependency":
			insertText = fmt.Sprintf("%s \"$name\" {\n\tconfig_path = $0\n}", blockName)
		case "generate", "include":
			insertText = fmt.Sprintf("%s \"$name\" {\n\tpath = $0\n}", blockName)
		case "terraform":
			insertText = fmt.Sprintf("%s {\n\tsource = $0\n}", blockName)
		default:
			insertText = fmt.Sprintf("%s {\n\t$0\n}", blockName)
		}

		item := lsp.CompletionItem{
			Kind:             lsp.CompletionItemKindSnippet,
			Label:            fmt.Sprintf("%s", blockName),
			Detail:           blockName,
			InsertText:       insertText,
			InsertTextFormat: lsp.InsertTextFormatSnippet,
			Documentation:    createMarkupContent(content),
		}

		*items = append(*items, item)
	}
}

func getAttributesTopLevel(items *[]lsp.CompletionItem) {
	docs, err := getDocumentationContentsBasedOnDirName("attributes-top-level")
	if err != nil {
		log.Panic("Failed to read arguments documentation")
		panic(1)
	}

	for _, doc := range docs {
		docLines := strings.Split(doc, "\n")

		attrName := docLines[0]
		content := strings.Join(docLines[1:], "\n")

		// A couple of attributes have a different type, so we build the
		// InsertText field based on attrName
		var insertText string
		switch attrName {
		case "inputs":
			insertText = fmt.Sprintf("%s = {\n\t$0\n}", attrName)
		default:
			insertText = fmt.Sprintf("%s = \"$0\"", attrName)
		}

		item := lsp.CompletionItem{
			Kind:             lsp.CompletionItemKindProperty,
			Label:            fmt.Sprintf("%s", attrName),
			Detail:           attrName,
			InsertText:       insertText,
			InsertTextFormat: lsp.InsertTextFormatSnippet,
			Documentation:    createMarkupContent(content),
		}

		*items = append(*items, item)
	}
}

func getAttributes(items *[]lsp.CompletionItem) {
	docs, err := getDocumentationContentsBasedOnDirName("attributes")
	if err != nil {
		log.Panic("Failed to read arguments documentation")
		panic(1)
	}

	for _, doc := range docs {
		docLines := strings.Split(doc, "\n")

		attrName := docLines[0]
		content := strings.Join(docLines[1:], "\n")

		// A couple of attributes have a different type, so we build the
		// InsertText field based on attrName
		var insertText string
		switch attrName {
		case "retryable_errors":
			insertText = fmt.Sprintf("%s = [$0]", attrName)
		case "skip", "prevent_destroy":
			insertText = fmt.Sprintf("%s = true", attrName)
		case "iam_assume_role_duration":
			insertText = fmt.Sprintf("%s = 14400", attrName)
		default:
			insertText = fmt.Sprintf("%s = \"$0\"", attrName)
		}

		item := lsp.CompletionItem{
			Kind:             lsp.CompletionItemKindProperty,
			Label:            fmt.Sprintf("%s", attrName),
			Detail:           attrName,
			InsertText:       insertText,
			InsertTextFormat: lsp.InsertTextFormatSnippet,
			Documentation:    createMarkupContent(content),
		}

		*items = append(*items, item)
	}
}
