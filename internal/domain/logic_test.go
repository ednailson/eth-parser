package domain

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ednailson/eth-parser/internal/domain/entity"
	"github.com/ednailson/eth-parser/pkg/blockchain"
	"github.com/ednailson/eth-parser/pkg/store"
)

func TestGetCurrentBlock(t *testing.T) {
	tests := []struct {
		name           string
		expectedOutput int
		blockchain     blockchain.Client
		store          store.Service
	}{
		{
			name:           "get current block successfully",
			expectedOutput: 1207,
			blockchain: &blockchainMock{
				block: "0x4b7",
				t:     t,
			},
			store: &storeMock{t: t},
		},
		{
			name:           "get current block with error",
			expectedOutput: -1,
			blockchain: &blockchainMock{
				block: "",
				err:   fmt.Errorf("error"),
				t:     t,
			},
			store: &storeMock{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logic := NewLogic(test.blockchain, test.store)

			output := logic.GetCurrentBlock()

			if output != test.expectedOutput {
				t.Errorf("output %d\nexpected output %d", output, test.expectedOutput)
			}
		})
	}
}

func TestSubscribe(t *testing.T) {
	tests := []struct {
		name           string
		expectedOutput bool
		input          string
		blockchain     blockchain.Client
		store          store.Service
	}{
		{
			name:           "subscribe successfully with address not registered",
			expectedOutput: true,
			input:          "eth_address_test",
			blockchain:     &blockchainMock{t: t},
			store: &storeMock{
				t:             t,
				validateInput: true,
				inputAddress:  "eth_address_test",
				err:           fmt.Errorf("no transactions error"),
			},
		},
		{
			name:           "subscribe successfully with address registered",
			expectedOutput: true,
			input:          "eth_address_test",
			blockchain: &blockchainMock{
				t:            t,
				transactions: []entity.Transaction{},
			},
			store: &storeMock{
				t: t,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logic := NewLogic(test.blockchain, test.store)

			output := logic.Subscribe(test.input)

			if output != test.expectedOutput {
				t.Errorf("output %v\nexpected output %v", output, test.expectedOutput)
			}
		})
	}
}

func TestGetTransactions(t *testing.T) {
	tests := []struct {
		name           string
		expectedOutput []entity.Transaction
		input          string
		blockchain     blockchain.Client
		store          store.Service
	}{
		{
			name: "get transactions successfully",
			expectedOutput: []entity.Transaction{
				{
					From: "eth_address_test",
					To:   "wallet_2_test",
				},
				{
					From: "wallet_3_test",
					To:   "eth_address_test",
				},
			},
			input:      "eth_address_test",
			blockchain: &blockchainMock{t: t},
			store: &storeMock{
				t:             t,
				validateInput: true,
				inputAddress:  "eth_address_test",
				transactions: []entity.Transaction{
					{
						From: "eth_address_test",
						To:   "wallet_2_test",
					},
					{
						From: "wallet_3_test",
						To:   "eth_address_test",
					},
				},
			},
		},
		{
			name:           "get transactions with error",
			expectedOutput: nil,
			input:          "eth_address_test",
			blockchain:     &blockchainMock{t: t},
			store: &storeMock{
				t:             t,
				validateInput: true,
				inputAddress:  "eth_address_test",
				err:           fmt.Errorf("failed to get transactions"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logic := NewLogic(test.blockchain, test.store)

			output := logic.GetTransactions(test.input)

			if !reflect.DeepEqual(output, test.expectedOutput) {
				t.Errorf("output %v\nexpected output %v", output, test.expectedOutput)
			}
		})
	}
}

type blockchainMock struct {
	err                          error
	block                        string
	transactions                 []entity.Transaction
	expectedGetTransactionsInput string
	validateInput                bool
	t                            *testing.T
}

func (b *blockchainMock) GetCurrentBlock() (string, error) {
	return b.block, b.err
}

func (b *blockchainMock) GetTransactions(block string) ([]entity.Transaction, error) {
	if b.validateInput {
		if block != b.expectedGetTransactionsInput {
			b.t.Errorf("input %s\nexpected input %s", block, b.expectedGetTransactionsInput)
		}
	}

	return b.transactions, b.err
}

type storeMock struct {
	t                 *testing.T
	err               error
	inputAddress      string
	validateInput     bool
	transactions      []entity.Transaction
	inputTransactions []entity.Transaction
	addresses         []string
}

func (s *storeMock) GetAddresses() []string {
	return s.addresses
}

func (s *storeMock) GetTransactions(address string) ([]entity.Transaction, error) {
	if s.validateInput {
		if address != s.inputAddress {
			s.t.Errorf("input %s\nexpected input %s", address, s.inputAddress)
		}
	}
	return s.transactions, s.err
}

func (s *storeMock) SaveTransactions(address string, transactions []entity.Transaction) {
	if s.validateInput {
		if s.inputAddress != address {
			s.t.Errorf("input %s\nexpected input %s", address, s.inputAddress)
		}

		if !reflect.DeepEqual(transactions, s.inputTransactions) {
			s.t.Errorf("input %v\nexpected input %v", transactions, s.inputTransactions)
		}
	}
}

func (s *storeMock) SaveAddress(address string) {}
