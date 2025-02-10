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

	// case lsp.MethodTextDocumentDidSave:
	// 	return h.didSave(ctx, reply, req)

	case lsp.MethodTextDocumentDidChange:
		log.Debug("Method for document change: ", req.Method())
		return h.didChange(ctx, reply, req)

	case lsp.MethodTextDocumentDidClose:
		log.Debug("Method for document close: ", req.Method())
		return h.didClose(ctx, reply, req)

	case lsp.MethodTextDocumentCompletion:
		log.Debug("Method for completion: ", req.Method())
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
				// TODO: check why incremental does not work
				// Change:    lsp.TextDocumentSyncKindIncremental,
				Change:    lsp.TextDocumentSyncKindFull,
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

func (h *handler) didChange(
	ctx context.Context,
	reply jsonrpc2.Replier,
	req jsonrpc2.Request,
) error {
	var params lsp.DidChangeTextDocumentParams
	json.Unmarshal(req.Params(), &params)
	h.documents.UpdateDocumentFull(params)
	// h.documents.UpdateDocumentIncremental(params)

	return reply(ctx, nil, nil)
}

func (h *handler) didClose(
	ctx context.Context,
	reply jsonrpc2.Replier,
	req jsonrpc2.Request,
) error {
	var params lsp.DidCloseTextDocumentParams
	json.Unmarshal(req.Params(), &params)
	h.documents.PopDocument(params)

	return reply(ctx, nil, nil)
}
