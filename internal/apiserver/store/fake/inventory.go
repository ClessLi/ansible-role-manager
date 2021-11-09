package fake

import (
	"context"
	ansible_inventory "github.com/ClessLi/ansible-role-manager/internal/pkg/ansible-inventory"
	metav1 "github.com/ClessLi/ansible-role-manager/internal/pkg/meta/v1"
)

type inventory struct {
	ds *datastore
}

func newInventory(ds *datastore) *inventory {
	return &inventory{ds: ds}
}

func (i inventory) Create(ctx context.Context, group ansible_inventory.Group, options metav1.CreateOptions) error {
	panic("implement me")
}

func (i inventory) Delete(ctx context.Context, groupName string, options metav1.DeleteOptions) error {
	panic("implement me")
}

func (i inventory) DeleteCollection(ctx context.Context, groupNames []string, options metav1.DeleteOptions) error {
	panic("implement me")
}

func (i inventory) Update(ctx context.Context, group ansible_inventory.Group, options metav1.UpdateOptions) error {
	panic("implement me")
}

func (i inventory) Get(ctx context.Context, groupName string, options metav1.GetOptions) (ansible_inventory.Group, error) {
	panic("implement me")
}

func (i inventory) List(ctx context.Context, options metav1.ListOptions) (*ansible_inventory.Groups, error) {
	panic("implement me")
}
