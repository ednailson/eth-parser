package domain

import "github.com/ednailson/eth-parser/internal/domain/entity"

type Parser interface {
	// GetCurrentBlock gets current block
	GetCurrentBlock() int
	// Subscribe adds an address to blockchain observer
	Subscribe(address string) bool
	// GetTransactions gets address' transactions after server boot and address' subscription
	GetTransactions(address string) []entity.Transaction
}
