package order

import (
	def "github.com/DeDevir/go_homework/order/internal/repository"
	"github.com/DeDevir/go_homework/order/internal/repository/model"
	"sync"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[string]*model.Order
}

func NewRepository() *repository {
	return &repository{
		data: make(map[string]*model.Order),
	}
}
