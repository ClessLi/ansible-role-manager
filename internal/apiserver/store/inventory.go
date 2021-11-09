package store

import (
	"context"
	ansible_inventory "github.com/ClessLi/ansible-role-manager/internal/pkg/ansible-inventory"
	metav1 "github.com/ClessLi/ansible-role-manager/internal/pkg/meta/v1"
)

type InventoryStore interface {
	Create(ctx context.Context, group ansible_inventory.Group, options metav1.CreateOptions) error
	Delete(ctx context.Context, groupName string, options metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, groupNames []string, options metav1.DeleteOptions) error
	Update(ctx context.Context, group ansible_inventory.Group, options metav1.UpdateOptions) error
	Get(ctx context.Context, groupName string, options metav1.GetOptions) (ansible_inventory.Group, error)
	List(ctx context.Context, options metav1.ListOptions) (*ansible_inventory.Groups, error)
}
