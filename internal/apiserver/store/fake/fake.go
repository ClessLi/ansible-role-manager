package fake

import (
	"fmt"
	"github.com/ClessLi/ansible-role-manager/internal/apiserver/store"
	ansible_inventory "github.com/ClessLi/ansible-role-manager/internal/pkg/ansible-inventory"
	"sync"
)

const ResourceCount = 255

type datastore struct {
	sync.RWMutex
	parser ansible_inventory.Parser
	inv    ansible_inventory.Inventory
}

func (ds *datastore) Inventory() store.InventoryStore {
	return newInventory(ds)
}

var (
	fakeFactory store.Factory
	once        sync.Once
)

func GetFakeFactoryOr() (store.Factory, error) {
	once.Do(func() {
		fakeFactory = &datastore{
			parser: ansible_inventory.NewParser(),
			inv:    ansible_inventory.NewInventory(FakeGroups(ResourceCount)),
		}
	})

	if fakeFactory == nil {
		return nil, fmt.Errorf("failed to get file store factory, fileFactory: %+v", fakeFactory)
	}

	return fakeFactory, nil
}

func FakeGroups(count int) map[string]ansible_inventory.Group {
	groups := make(map[string]ansible_inventory.Group)

	for i := 0; i < count; i++ {
		groupname := fmt.Sprintf("test-group%d", i)
		groups[groupname], _ = ansible_inventory.NewGroup(groupname, []ansible_inventory.Host{
			ansible_inventory.ParseHost(fmt.Sprintf("192.168.%d.1", i)),
			ansible_inventory.ParseHost(fmt.Sprintf("10.%d.0.1", i)),
			ansible_inventory.ParseHost(fmt.Sprintf("192.168.%d.[%d:254]", i, i/2)),
		})
	}

	return groups
}
