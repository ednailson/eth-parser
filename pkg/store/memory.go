package store

import (
	"fmt"
	"slices"
	"sync"

	"github.com/ednailson/eth-parser/internal/domain/entity"
)

type Memory struct {
	addresses    []string
	transactions map[string][]entity.Transaction
	mutex        sync.RWMutex
}

func NewMemory() *Memory {
	return &Memory{
		transactions: make(map[string][]entity.Transaction),
		mutex:        sync.RWMutex{},
	}
}

func (m *Memory) GetAddresses() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.addresses
}

func (m *Memory) GetTransactions(address string) ([]entity.Transaction, error) {
	m.mutex.RLock()
	transactions, exists := m.transactions[address]
	m.mutex.RUnlock()
	if !exists {
		return nil, fmt.Errorf("no transaction found")
	}

	return transactions, nil
}

func (m *Memory) SaveTransactions(address string, transactions []entity.Transaction) {
	m.mutex.Lock()
	m.transactions[address] = append(m.transactions[address], transactions...)
	m.mutex.Unlock()

	m.SaveAddress(address)
}

func (m *Memory) SaveAddress(address string) {
	if !slices.Contains(m.addresses, address) {
		m.mutex.Lock()
		m.addresses = append(m.addresses, address)
		m.mutex.Unlock()
	}
}
