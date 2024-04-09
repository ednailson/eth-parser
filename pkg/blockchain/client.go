package blockchain

import "github.com/ednailson/eth-parser/internal/domain/entity"

type Client interface {
	// GetCurrentBlock fetches current blockchain block.
	GetCurrentBlock() (string, error)
	// GetTransactions fetches transactions of given address.
	GetTransactions(block string) ([]entity.Transaction, error)
}
