package file

import (
	"fmt"
	"github.com/ClessLi/ansible-role-manager/internal/apiserver/store"
	"github.com/ClessLi/ansible-role-manager/internal/pkg/ansible-inventory"
	"github.com/ClessLi/ansible-role-manager/internal/pkg/file-store"
	"sync"
)

var (
	fileFactory store.Factory
	once        = sync.Once{}
)

type fileStore struct {
	fs file_store.FileStore
}

func (f *fileStore) Inventory() store.InventoryStore {
	return newInventory(f, ansible_inventory.NewParser())
}

func GetFileFactory() (store.Factory, error) {
	var err error
	var storeIns file_store.FileStore
	once.Do(func() {
		if fileFactory == nil {
			//dirPath := filepath.Join(config.ExtConfig.Ansible.BaseDir, config.ExtConfig.Ansible.InventoryDir)
			// TODO: 加载配置文件，并解析出inventory目录
			//parser := ansible_inventory.NewParser()
			//storeIns, err = newFileStore(dirPath, parser)
			storeIns, err = file_store.NewFileStore(dirPath)
			fileFactory = &fileStore{fs: storeIns}
		}
	})

	if fileFactory == nil || err != nil {
		return nil, fmt.Errorf("failed to get file store fatory, fileFactory: %+v, error: %w", fileFactory, err)
	}

	return fileFactory, nil
}
