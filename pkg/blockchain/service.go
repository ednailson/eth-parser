package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"

	"github.com/ednailson/eth-parser/internal/domain/entity"
)

type Service struct {
	httpClient *http.Client
	host       string
}

func NewService(host string) *Service {
	return &Service{
		httpClient: http.DefaultClient,
		host:       host,
	}
}

func (s *Service) GetCurrentBlock() (string, error) {
	request := entity.JSONRPCRequestBytes(1, blockNumberMethod, nil)

	var result string
	if err := s.doRequest(request, &result); err != nil {
		return "", fmt.Errorf("failed to do get current block request, error %w", err)
	}

	return result, nil
}

func (s *Service) GetTransactions(block string) ([]entity.Transaction, error) {
	var result []entity.Transaction

	transactions, err := s.getBlockTransactions(block)
	if err != nil {
		return nil, err
	}

	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}

	for _, transaction := range transactions {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			params, err := json.Marshal([]string{t})
			if err != nil {
				slog.Error("Failed to marshal get transaction parameters", "error", err)
				return
			}
			request := entity.JSONRPCRequestBytes(2, getTransactionMethod, params)

			var res entity.Transaction
			if err := s.doRequest(request, &res); err != nil {
				slog.Error("Failed to do get transactions request", "error", err)
				return
			}

			mutex.Lock()
			result = append(result, res)
			mutex.Unlock()
		}(transaction)
	}

	wg.Wait()

	return result, nil
}

func (s *Service) getBlockTransactions(block string) ([]string, error) {
	params, err := json.Marshal([]any{block, false})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal get block transactions parameters, error: %w", err)
	}
	request := entity.JSONRPCRequestBytes(3, getBlockInfoMethod, params)

	var result GetBlockTransactions
	if err := s.doRequest(request, &result); err != nil {
		return nil, fmt.Errorf("failed to do get block transactions request, error: %w", err)
	}

	return result.Transactions, nil
}

func (s *Service) doRequest(requestBody []byte, data any) error {
	response, err := s.httpClient.Post(s.host, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to request, error: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read body from request, error: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to request, status not ok, status received %d, body %s", response.StatusCode, string(body))
	}

	var resp entity.JSONRPCResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return fmt.Errorf("failed to unmarshal body from response, result %v, error: %w", string(body), err)
	}

	if resp.Error != nil {
		return fmt.Errorf("response has error, error: %v", resp.Error)
	}

	err = json.Unmarshal(resp.Result, &data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal result from get transactions response, result %v, error: %w", string(resp.Result), err)
	}

	return nil
}
