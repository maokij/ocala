package langserver

import (
	"context"
	"encoding/json"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *langHandler) handleTextDocumentDidChange(_ context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params DidChangeTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	changes := params.ContentChanges
	if len(changes) == 1 && changes[0].Range == nil && changes[0].RangeLength == 0 {
		if err := h.updateFile(params.TextDocument.URI, changes[0].Text, &params.TextDocument.Version, eventTypeChange); err != nil {
			return nil, err
		}
	} else {
		if err := h.updateFileIncremental(params.TextDocument.URI, params.ContentChanges, &params.TextDocument.Version, eventTypeChange); err != nil {
			return nil, err
		}
	}
	return nil, nil
}
