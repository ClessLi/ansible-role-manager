package v1

import (
	"context"
	"github.com/ClessLi/ansible-role-manager/internal/apiserver/store"
	ansible_inventory "github.com/ClessLi/ansible-role-manager/internal/pkg/ansible-inventory"
	metav1 "github.com/ClessLi/ansible-role-manager/internal/pkg/meta/v1"
	"github.com/marmotedu/errors"
)

//type InventorySrv interface {
//	AddHostToGroup(groupName string, hosts ...ansible_inventory.Host) error
//	RenewGroupName(oldName, newName string) error
//	RemoveHostFromGroup(groupName string, hosts ...ansible_inventory.Host) error
//	RemoveGroup(groupName string) error
//	GetGroupsByPage(limit, page int) *ansible_inventory.Groups
//}

type InventorySrv interface {
	Create(ctx context.Context, group ansible_inventory.Group, options metav1.CreateOptions) error
	Delete(ctx context.Context, groupName string, options metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, groupNames []string, options metav1.DeleteOptions) error
	Update(ctx context.Context, group ansible_inventory.Group, options metav1.UpdateOptions) error
	Get(ctx context.Context, groupName string, options metav1.GetOptions) (ansible_inventory.Group, error)
	List(ctx context.Context, options metav1.ListOptions) (*ansible_inventory.Groups, error)
}

type inventoryService struct {
	store store.Factory
}

var _ InventorySrv = (*inventoryService)(nil)

func newInventory(srv *service) *inventoryService {
	return &inventoryService{store: srv.store}
}

func (i *inventoryService) Create(ctx context.Context, group ansible_inventory.Group, options metav1.CreateOptions) error {
	return errors.Wrapf(i.store.Inventory().Create(ctx, group, options), "creat group '%v' error", group.GetName())
}

func (i *inventoryService) Delete(ctx context.Context, groupName string, options metav1.DeleteOptions) error {
	return errors.Wrapf(i.store.Inventory().Delete(ctx, groupName, options), "delete group '%v' error", groupName)
}

func (i *inventoryService) DeleteCollection(ctx context.Context, groupNames []string, options metav1.DeleteOptions) error {
	return errors.Wrapf(i.store.Inventory().DeleteCollection(ctx, groupNames, options), "delete groups '%v' error", groupNames)
}

func (i *inventoryService) Update(ctx context.Context, group ansible_inventory.Group, options metav1.UpdateOptions) error {
	return errors.Wrapf(i.store.Inventory().Update(ctx, group, options), "update group '%v' error", group.GetName())
}

func (i *inventoryService) Get(ctx context.Context, groupName string, options metav1.GetOptions) (ansible_inventory.Group, error) {
	group, err := i.store.Inventory().Get(ctx, groupName, options)
	return group, errors.Wrapf(err, "get group '%v' error", groupName)
}

func (i *inventoryService) List(ctx context.Context, options metav1.ListOptions) (*ansible_inventory.Groups, error) {
	groups, err := i.store.Inventory().List(ctx, options)
	return groups, errors.Wrapf(err, "list groups error")
}
