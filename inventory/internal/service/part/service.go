package part

import "github.com/DeDevir/go_homework/inventory/internal/repository"

type service struct {
	partRepository repository.PartRepository
}

func NewService(repository repository.PartRepository) *service {
	return &service{
		partRepository: repository,
	}
}
