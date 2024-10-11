package handler

import (
	"context"
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"go.lsp.dev/jsonrpc2"
	lsp "go.lsp.dev/protocol"

	_ "github.com/ArnauLlamas/terragrunt-ls/internal/log"
	lsplocal "github.com/ArnauLlamas/terragrunt-ls/internal/lsp"
)

type handler struct {
	connPool  jsonrpc2.Conn
	documents *lsplocal.DocumentStore
}

func NewHandler(connPool jsonrpc2.Conn) jsonrpc2.Handler {
	documents, _ := lsplocal.NewDocumentStore("")
	handler := &handler{
		connPool:  connPool,
		documents: documents,
	}
	return jsonrpc2.ReplyHandler(handler.reqHandler)
}

func (h *handler) reqHandler(
	ctx context.Context,
	reply jsonrpc2.Replier,
	req jsonrpc2.Request,
) error {
	log.Debug("request received: ", req)

	switch req.Method() {
	case lsp.MethodInitialize:
		return h.initialize(ctx, reply, req)

	case lsp.MethodInitialized:
		return reply(ctx, nil, nil)

	case lsp.MethodShutdown:
		return h.shutdown()

	case lsp.MethodTextDocumentDidOpen:
		return h.didOpen(ctx, reply, req)

	case lsp.MethodTextDocumentCompletion:
		return h.completion(ctx, reply, req)

	default:
		log.Info("Method not supported", req.Method())
	}

	return jsonrpc2.MethodNotFoundHandler(ctx, reply, req)
}

func (h *handler) initialize(
	ctx context.Context,
	reply jsonrpc2.Replier,
	_ jsonrpc2.Request,
) error {
	return reply(ctx, lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync: lsp.TextDocumentSyncOptions{
				// Change:    lsp.TextDocumentSyncKindIncremental,
				OpenClose: true,
				Save:      &lsp.SaveOptions{},
			},
			CompletionProvider: &lsp.CompletionOptions{
				ResolveProvider: false,
			},
		},
	}, nil)
}

func (h *handler) shutdown() error {
	return h.connPool.Close()
}

func (h *handler) didOpen(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params lsp.DidOpenTextDocumentParams
	json.Unmarshal(req.Params(), &params)
	h.documents.PushDocument(params)

	return reply(ctx, nil, nil)
}
