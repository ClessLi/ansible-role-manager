package file

import (
	"context"
	"errors"
	"fmt"
	ansible_inventory "github.com/ClessLi/ansible-role-manager/internal/pkg/ansible-inventory"
	file_store "github.com/ClessLi/ansible-role-manager/internal/pkg/file-store"
	metav1 "github.com/ClessLi/ansible-role-manager/internal/pkg/meta/v1"
)

type inventory struct {
	fs     file_store.FileStore
	parser ansible_inventory.Parser
}

func (i *inventory) Create(ctx context.Context, group ansible_inventory.Group, options metav1.CreateOptions) error {
	inv, err := i.load()
	if err != nil {
		return err
	}

	if inv.GetAllGroups()[group.GetName()] != nil {
		return fmt.Errorf("group '%v' is exist", group.GetName())
	}

	err = inv.AddHostToGroup(group.GetName(), group.GetHosts()...)
	if err != nil {
		return err
	}

	return i.save(inv)
}

func (i *inventory) Delete(ctx context.Context, groupName string, options metav1.DeleteOptions) error {
	// 加载inventory
	inv, err := i.load()
	if err != nil {
		return err
	}

	// 检查是否存在该主机组
	if inv.GetAllGroups()[groupName] == nil && !options.Force {
		return fmt.Errorf("group '%v' is not exist", groupName)
	}

	// 删除该主机组
	inv.RemoveGroup(groupName)

	// 保存inventory
	return i.save(inv)
}

func (i *inventory) DeleteCollection(ctx context.Context, groupNames []string, options metav1.DeleteOptions) error {
	// 加载inventory
	inv, err := i.load()
	if err != nil {
		return err
	}

	// 是否非强制操作
	if !options.Force {
		// 非强制操作则开始校验传参
		if groupNames == nil || len(groupNames) < 1 {
			return errors.New("the 'groupNames' param is nil or null")
		}

		nonexistentGroups := make([]string, 0)

		for _, groupName := range groupNames {
			if inv.GetAllGroups()[groupName] == nil {
				nonexistentGroups = append(nonexistentGroups, groupName)
			}
		}

		if len(nonexistentGroups) > 0 {
			return fmt.Errorf("nonexistent group list: %+v", nonexistentGroups)
		}
	}

	// groupNames为nil或空时跳过处理
	if groupNames == nil || len(groupNames) < 1 {
		return nil
	}

	// 删除该主机组列表
	for _, groupName := range groupNames {
		inv.RemoveGroup(groupName)
	}

	// 保存inventory
	return i.save(inv)
}

func (i *inventory) Update(ctx context.Context, group ansible_inventory.Group, options metav1.UpdateOptions) error {
	// 加载inventory
	inv, err := i.load()
	if err != nil {
		return err
	}

	// 检查是否存在该主机组
	if inv.GetAllGroups()[group.GetName()] == nil {
		return fmt.Errorf("group '%v' is not exist", group.GetName())
	}

	// 删除该主机组
	inv.RemoveGroup(group.GetName())

	// 添加组到inventory
	err = inv.AddHostToGroup(group.GetName(), group.GetHosts()...)
	if err != nil {
		return err
	}

	// 保存inventory
	return i.save(inv)
}

func (i *inventory) Get(ctx context.Context, groupName string, options metav1.GetOptions) (ansible_inventory.Group, error) {
	// 加载inventory
	inv, err := i.load()
	if err != nil {
		return nil, err
	}

	group := inv.GetAllGroups()[groupName]

	// 检查是否存在该主机组
	if group == nil {
		return nil, fmt.Errorf("group '%v' is not exist", groupName)
	}

	return group, nil
}

func (i *inventory) List(ctx context.Context, options metav1.ListOptions) (*ansible_inventory.Groups, error) {
	// 加载inventory
	inv, err := i.load()
	if err != nil {
		return nil, err
	}

	return inv.GetGroupsByPage(*options.Page, *options.NumPerPage)
}

func (i *inventory) load() (ansible_inventory.Inventory, error) {
	filePaths, err := i.fs.AllFiles()
	if err != nil {
		return nil, err
	}
	groups := make(map[string]ansible_inventory.Group)
	for _, filePath := range filePaths {
		b, err := i.fs.Read(filePath)
		if err != nil {
			return nil, err
		}
		g, err := i.parser.Parse(b)
		if err != nil {
			fmt.Printf("parse inventory file %s failed, cased by: %s\n", filePath, err)
			continue
		}
		if _, has := groups[g.GetName()]; !has {
			groups[g.GetName()] = g
			continue
		}
		_ = groups[g.GetName()].AddHost(g.GetHosts()...)
	}

	inv := ansible_inventory.NewInventory(groups)
	return inv, nil
}

func (i *inventory) save(inv ansible_inventory.Inventory) error {
	// 清理TruncatedGroup
	for tGName, isTruncated := range inv.GetTruncatedGroup() {
		if isTruncated {
			//err = os.Remove(filepath.Join(i.dirPath, tGName))
			err := i.fs.Remove(tGName)
			if err != nil {
				return err
			}
		}
	}

	// 保存配置
	for gName, g := range inv.GetAllGroups() {
		b, err := i.parser.Dump(g)
		if err != nil {
			return err
		}

		//err = os.WriteFile(filepath.Join(*i.dirPath, gName), b, 0644)
		err = i.fs.Write(gName, b)
		if err != nil {
			return err
		}
	}
	return nil
}

func newInventory(fs *fileStore, parser ansible_inventory.Parser) *inventory {
	return &inventory{
		fs:     fs.fs,
		parser: parser,
	}
}
