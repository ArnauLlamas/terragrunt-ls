package lsp

import (
	"net/url"
	"os"

	sitter "github.com/smacker/go-tree-sitter"
	lsp "go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

// documentStore holds opened documents.
type DocumentStore struct {
	documents  map[string]*Document
	workingDir string
}

func NewDocumentStore(workingDir string) (*DocumentStore, error) {
	if workingDir == "" {
		_, err := os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	return &DocumentStore{
		documents:  map[string]*Document{},
		workingDir: workingDir,
	}, nil
}

func (s *DocumentStore) GetDocument(docuri uri.URI) *Document {
	path := getPathFromURI(docuri)
	return s.documents[path]
}

func (s *DocumentStore) PushDocument(params lsp.DidOpenTextDocumentParams) *Document {
	path := getPathFromURI(params.TextDocument.URI)
	doc := &Document{
		URI:     params.TextDocument.URI,
		Path:    path,
		Content: params.TextDocument.Text,
		Ast:     ParseAst([]byte(params.TextDocument.Text)),
	}

	s.documents[path] = doc
	return doc
}

func getPathFromURI(docuri uri.URI) string {
	parsed, _ := url.Parse(docuri.Filename())
	return parsed.Path
}

// Document represents an opened file.
type Document struct {
	URI     lsp.DocumentURI
	Path    string
	Content string
	lines   []string
	Ast     *sitter.Tree
}
