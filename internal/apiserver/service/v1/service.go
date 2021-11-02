package v1

import "github.com/ClessLi/ansible-role-manager/internal/apiserver/store"

type Service interface {
	Inventory() InventorySrv
}

type service struct {
	store store.Factory
}

func NewService(store store.Factory) Service {
	return Service(&service{store: store})
}

func (s *service) Inventory() InventorySrv {
	return newInventory(s)
}
