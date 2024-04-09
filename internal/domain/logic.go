package domain

import (
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/ednailson/eth-parser/internal/domain/entity"
	"github.com/ednailson/eth-parser/pkg/blockchain"
	"github.com/ednailson/eth-parser/pkg/store"
)

type Logic struct {
	blockchain         blockchain.Client
	store              store.Service
	currentBlock       string
	latestTransactions []entity.Transaction
	mutex              sync.RWMutex
}

func NewLogic(b blockchain.Client, s store.Service) *Logic {
	logic := &Logic{
		blockchain: b,
		store:      s,
	}

	logic.observer()

	go func(l *Logic) {
		// Ticker duration should be lowered, but it is one minute because there is a rate limit in cloudflare.
		ticker := time.NewTicker(time.Minute)
		for range ticker.C {
			l.observer()
		}
	}(logic)

	return logic
}

func (l *Logic) observer() {
	block, err := l.blockchain.GetCurrentBlock()
	if err != nil {
		slog.Error("Failed to get current block from blockchain", "error", err)
		return
	}

	l.mutex.RLock()
	isSameBlock := block == l.currentBlock
	l.mutex.RUnlock()
	if isSameBlock {
		return
	}

	transactions, err := l.blockchain.GetTransactions(block)
	if err != nil {
		slog.Error("Failed to get transactions from blockchain", "block", block, "error", err)
		return
	}

	l.mutex.Lock()
	l.currentBlock = block
	l.latestTransactions = transactions
	l.mutex.Unlock()

	for _, addr := range l.store.GetAddresses() {
		addressTransactions := l.addressTransactions(addr)

		l.store.SaveTransactions(addr, addressTransactions)
	}
}

func (l *Logic) GetCurrentBlock() int {
	l.mutex.RLock()
	block := l.currentBlock
	l.mutex.RUnlock()

	number, err := hexToInt(block)
	if err != nil {
		slog.Error("Failed to transform current block into int", "error", err)
		return -1
	}

	return number
}

func (l *Logic) Subscribe(address string) bool {
	_, err := l.store.GetTransactions(address)
	if err == nil {
		// Address is already being observed.
		return true
	}

	addressTransactions := l.addressTransactions(address)

	l.store.SaveTransactions(address, addressTransactions)
	return true
}

func (l *Logic) GetTransactions(address string) []entity.Transaction {
	transactions, err := l.store.GetTransactions(address)
	if err != nil {
		slog.Error("Failed to get transaction from store", "error", err)
		return nil
	}

	return transactions
}

func (l *Logic) addressTransactions(address string) []entity.Transaction {
	var addressTransactions []entity.Transaction
	l.mutex.RLock()
	for _, transaction := range l.latestTransactions {
		if transaction.From == address || transaction.To == address {
			addressTransactions = append(addressTransactions, transaction)
		}
	}
	l.mutex.RUnlock()
	return addressTransactions
}

func hexToInt(value string) (int, error) {
	number, err := strconv.ParseInt(value, 0, 64)
	if err != nil {
		return -1, fmt.Errorf("failed to parse block to int, error %w", err)
	}

	return int(number), nil
}
