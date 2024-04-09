package store

import "github.com/ednailson/eth-parser/internal/domain/entity"

type Service interface {
	GetAddresses() []string
	GetTransactions(address string) ([]entity.Transaction, error)
	SaveTransactions(address string, transaction []entity.Transaction)
	SaveAddress(address string)
}
