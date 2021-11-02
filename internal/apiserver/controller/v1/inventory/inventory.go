package inventory

import (
	srvv1 "github.com/ClessLi/ansible-role-manager/internal/apiserver/service/v1"
	"github.com/ClessLi/ansible-role-manager/internal/apiserver/store"
)

type InventoryController struct {
	srv srvv1.Service
}

func NewInventoryController(store store.Factory) *InventoryController {
	return &InventoryController{srv: srvv1.NewService(store)}
}
