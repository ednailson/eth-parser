package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ednailson/eth-parser/internal/domain"
)

type Server struct {
	logic domain.Parser
	port  int
}

func NewServer(port int, logic domain.Parser) *Server {
	return &Server{
		logic: logic,
		port:  port,
	}
}

func (s *Server) register() {
	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		req, failed := s.extractRequestBody(w, r)
		if failed {
			return
		}

		result := s.logic.Subscribe(req.Address)

		response, _ := json.Marshal(map[string]any{
			"result": result,
		})

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(response)
	})

	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		address := r.URL.Query().Get("address")
		if address == "" {
			w.WriteHeader(http.StatusBadRequest)
			response, _ := json.Marshal(map[string]any{
				"message": "address parameter is invalid",
			})
			_, _ = w.Write(response)
			return
		}

		transactions := s.logic.GetTransactions(address)

		response, _ := json.Marshal(map[string]any{
			"transactions": transactions,
		})

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(response)
	})

	http.HandleFunc("/currentBlock", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		block := s.logic.GetCurrentBlock()

		response, _ := json.Marshal(map[string]any{
			"current_block": block,
		})

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(response)
	})
}

func (s *Server) extractRequestBody(w http.ResponseWriter, r *http.Request) (*Request, bool) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response, _ := json.Marshal(map[string]any{
			"message": "failed to decode body",
		})
		_, _ = w.Write(response)
		return nil, true
	}

	var req Request
	err = json.Unmarshal(body, &req)
	if err != nil || req.Address == "" {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(map[string]any{
			"message": "address parameter is invalid",
		})
		_, _ = w.Write(response)
		return nil, true
	}

	return &req, false
}

func (s *Server) Serve() error {
	s.register()
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}
