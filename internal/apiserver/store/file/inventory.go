package file

import (
	"context"
	ansible_inventory "github.com/ClessLi/ansible-role-manager/internal/pkg/ansible-inventory"
	"github.com/ClessLi/ansible-role-manager/internal/pkg/code"
	file_store "github.com/ClessLi/ansible-role-manager/internal/pkg/file-store"
	metav1 "github.com/ClessLi/ansible-role-manager/internal/pkg/meta/v1"
	log "github.com/ClessLi/ansible-role-manager/pkg/log/v2"
	"github.com/marmotedu/errors"
)

type inventory struct {
	fs     file_store.FileStore
	parser ansible_inventory.Parser
}

func (i *inventory) Create(ctx context.Context, group ansible_inventory.Group, options metav1.CreateOptions) error {
	inv, err := i.load(ctx)
	if err != nil {
		return errors.WrapC(err, code.ErrDataRepository, "failed to load inventory")
	}

	if inv.GetAllGroups()[group.GetName()] != nil {
		return errors.WithCode(code.ErrGroupAlreadyExist, "group '%v' is exist", group.GetName())
	}

	err = inv.AddHostToGroup(group.GetName(), group.GetHosts()...)
	if err != nil {
		return errors.WithCode(code.ErrDataRepository, err.Error())
	}

	return errors.WrapC(i.save(ctx, inv), code.ErrDataRepository, "save inventory error")
}

func (i *inventory) Delete(ctx context.Context, groupName string, options metav1.DeleteOptions) error {
	// 加载inventory
	inv, err := i.load(ctx)
	if err != nil {
		return errors.WrapC(err, code.ErrDataRepository, "failed to load inventory")
	}

	// 检查是否存在该主机组
	if inv.GetAllGroups()[groupName] == nil && !options.Force {
		return errors.WithCode(code.ErrGroupNotFound, "group '%v' is not exist", groupName)
	}

	// 删除该主机组
	inv.RemoveGroup(groupName)

	// 保存inventory
	return errors.WrapC(i.save(ctx, inv), code.ErrDataRepository, "save inventory error")
}

func (i *inventory) DeleteCollection(ctx context.Context, groupNames []string, options metav1.DeleteOptions) error {
	// 加载inventory
	inv, err := i.load(ctx)
	if err != nil {
		return errors.WrapC(err, code.ErrDataRepository, "failed to load inventory")
	}

	// 是否非强制操作
	if !options.Force {
		// 非强制操作则开始校验传参
		if groupNames == nil || len(groupNames) < 1 {
			return errors.WithCode(code.ErrGroupNotFound, "the 'groupNames' param is nil or null")
		}

		nonexistentGroups := make([]string, 0)

		for _, groupName := range groupNames {
			if inv.GetAllGroups()[groupName] == nil {
				nonexistentGroups = append(nonexistentGroups, groupName)
			}
		}

		if len(nonexistentGroups) > 0 {
			return errors.WithCode(code.ErrGroupNotFound, "nonexistent group list: %v", nonexistentGroups)
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
	return errors.WrapC(i.save(ctx, inv), code.ErrDataRepository, "save inventory error")
}

func (i *inventory) Update(ctx context.Context, group ansible_inventory.Group, options metav1.UpdateOptions) error {
	// 加载inventory
	inv, err := i.load(ctx)
	if err != nil {
		return errors.WrapC(err, code.ErrDataRepository, "failed to load inventory")
	}

	// 检查是否存在该主机组
	if inv.GetAllGroups()[group.GetName()] == nil {
		return errors.WithCode(code.ErrGroupNotFound, "group '%v' is not exist", group.GetName())
	}

	// 删除该主机组
	inv.RemoveGroup(group.GetName())

	// 添加组到inventory
	err = inv.AddHostToGroup(group.GetName(), group.GetHosts()...)
	if err != nil {
		return errors.WithCode(code.ErrDataRepository, err.Error())
	}

	// 保存inventory
	return errors.WrapC(i.save(ctx, inv), code.ErrDataRepository, "save inventory error")
}

func (i *inventory) Get(ctx context.Context, groupName string, options metav1.GetOptions) (ansible_inventory.Group, error) {
	// 加载inventory
	inv, err := i.load(ctx)
	if err != nil {
		return nil, errors.WrapC(err, code.ErrDataRepository, "failed to load inventory")
	}

	group := inv.GetAllGroups()[groupName]

	// 检查是否存在该主机组
	if group == nil {
		return nil, errors.WithCode(code.ErrGroupNotFound, "group '%v' is not exist", groupName)
	}

	return group, nil
}

func (i *inventory) List(ctx context.Context, options metav1.ListOptions) (*ansible_inventory.Groups, error) {
	// 加载inventory
	inv, err := i.load(ctx)
	if err != nil {
		return nil, errors.WrapC(err, code.ErrDataRepository, "failed to load inventory")
	}

	groups, err := inv.GetGroupsByPage(*options.Page, *options.NumPerPage)
	return groups, errors.WrapC(err, code.ErrDataRepository, "list groups error")
}

func (i *inventory) load(ctx context.Context) (ansible_inventory.Inventory, error) {
	filePaths, err := i.fs.AllFiles()
	if err != nil {
		log.L(ctx).Warnf("fileStore get file paths failed %s", err.Error())
		return nil, errors.Wrap(err, "failed to get file paths")
	}
	groups := make(map[string]ansible_inventory.Group)
	errs := make([]error, 0)
	for _, filePath := range filePaths {
		b, err := i.fs.Read(filePath)
		if err != nil {
			log.L(ctx).Warnf("fileStore read file '%s' failed %s", filePath, err.Error())
			return nil, errors.Wrapf(err, "failed to read file '%s'", filePath)
		}
		g, err := i.parser.Parse(b)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "failed to parse inventory file %s", filePath))
			continue
		}
		if _, has := groups[g.GetName()]; !has {
			groups[g.GetName()] = g
			continue
		}
		_ = groups[g.GetName()].AddHost(g.GetHosts()...)
	}

	inv := ansible_inventory.NewInventory(groups)
	return inv, errors.NewAggregate(errs)
}

func (i *inventory) save(ctx context.Context, inv ansible_inventory.Inventory) error {
	// 清理TruncatedGroup
	for tGName, isTruncated := range inv.GetTruncatedGroup() {
		if isTruncated {
			//err = os.Remove(filepath.Join(i.dirPath, tGName))
			err := i.fs.Remove(tGName)
			if err != nil {
				log.L(ctx).Errorf("fileStore remove file '%s' failed %s", tGName, err.Error())
				return errors.New(err.Error())
			}
		}
	}

	// 保存配置
	for gName, g := range inv.GetAllGroups() {
		b, err := i.parser.Dump(g)
		if err != nil {
			log.L(ctx).Warnf("ansible_inventory.Parser dump group '%s' failed %s", gName, err.Error())
			return errors.New(err.Error())
		}

		//err = os.WriteFile(filepath.Join(*i.dirPath, gName), b, 0644)
		err = i.fs.Write(gName, b)
		if err != nil {
			log.L(ctx).Warnf("fileStore write file '%s' failed %s", gName, err.Error())
			return errors.New(err.Error())
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
