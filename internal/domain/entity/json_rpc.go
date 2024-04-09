package entity

import (
	"encoding/json"
	"log/slog"
)

const JSONRPCVersion = "2.0"

type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      any             `json:"id"`
}

type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *JSONRPCError   `json:"error,omitempty"`
	ID      json.RawMessage `json:"id"`
}

type JSONRPCError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func JSONRPCRequestBytes(id any, method string, params json.RawMessage) []byte {
	result, err := json.Marshal(JSONRPCRequest{
		JSONRPC: JSONRPCVersion,
		Method:  method,
		Params:  params,
		ID:      id,
	})
	if err != nil {
		slog.Error("Failed to marshal json rpc request", "error", err)
		return nil
	}

	return result
}
