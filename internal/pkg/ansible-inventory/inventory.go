package ansible_inventory

import (
	"fmt"
	"sort"
	"strings"
)

type Inventory interface {
	GetAllGroups() map[string]Group
	GetTruncatedGroup() map[string]bool
	GenerateGroupAll() Group
	AddHostToGroup(groupName string, hosts ...Host) error
	RenewGroupName(oldName, newName string) error
	RemoveHostFromGroup(groupName string, hosts ...Host)
	RemoveGroup(groupName string)
	GetGroupsByPage(page, recordsPerPage uint) (*Groups, error) // DONE: 分页查询机制、反馈主机总数
}

type inventory struct {
	sortedGroupNames []string
	groups           map[string]Group
	isTruncatedGroup map[string]bool
}

func NewInventory(groups map[string]Group) Inventory {
	inv := &inventory{
		sortedGroupNames: make([]string, 0),
		groups:           groups,
		isTruncatedGroup: make(map[string]bool),
	}
	for s := range groups {
		inv.sortedGroupNames = append(inv.sortedGroupNames, s)
	}
	inv.sortGroups()
	return Inventory(inv)
}

func (i *inventory) AddHostToGroup(groupName string, hosts ...Host) error {
	// 不允许新增"all"组名的组
	if strings.EqualFold(strings.ToLower(groupName), "all") {
		return fmt.Errorf("can not build a group which named '%s'", groupName)
	}
	var g Group
	if _, has := i.groups[groupName]; has {
		g = i.groups[groupName]
	} else {
		g = newGroup()
		err := g.setName(groupName)
		if err != nil {
			return err
		}
	}
	err := g.AddHost(hosts...)
	if err != nil {
		return err
	}
	i.groups[groupName] = g
	i.isTruncatedGroup[groupName] = false
	i.sortedGroupNames = append(i.sortedGroupNames, groupName)
	i.sortGroups()
	return nil
}

func (i *inventory) RenewGroupName(oldName, newName string) error {
	// 不允许新增"all"组名的组
	if strings.EqualFold(strings.ToLower(newName), "all") {
		return fmt.Errorf("rename '%s' is not allowed", newName)
	}
	if _, has := i.groups[oldName]; !has {
		return fmt.Errorf("nonexistent group by name %s", oldName)
	}
	if _, has := i.groups[newName]; has {
		return fmt.Errorf("duplicate group name %s", newName)
	}
	g := i.groups[oldName]
	i.RemoveGroup(oldName)
	return i.AddHostToGroup(newName, g.GetHosts()...)
}

func (i *inventory) RemoveHostFromGroup(groupName string, hosts ...Host) {
	if _, has := i.groups[groupName]; has {
		i.groups[groupName].RemoveHost(hosts...)
	}
}

func (i *inventory) RemoveGroup(groupName string) {
	if _, has := i.groups[groupName]; has {
		delete(i.groups, groupName)
		idx := i.searchGroup(groupName)
		i.sortedGroupNames = append(i.sortedGroupNames[:idx], i.sortedGroupNames[idx+1:]...)
		i.isTruncatedGroup[groupName] = true
	}
}

func (i inventory) GetGroupsByPage(page, recordsPerPage uint) (*Groups, error) {
	totalGroupsNum := uint(len(i.sortedGroupNames))
	var totalPagesNum uint
	groups := &Groups{GroupsMap: make(map[string]Group), TotalGroupsNum: &totalGroupsNum, TotalPagesNum: &totalPagesNum}

	if recordsPerPage == 0 {
		return groups, nil
	}

	// 默认0转换为1
	if page == 0 {
		page = 1
	}

	// 加入all组进行页数计算
	totalPagesNum = (totalGroupsNum + 1) / recordsPerPage
	if (totalGroupsNum+1)%recordsPerPage > 0 {
		totalPagesNum++
	}

	// 计算起始索引
	startIdx := (int(page) - 1) * int(recordsPerPage)
	// 因加入了all组进行计算，起始索引需往前移动一位
	startIdx--
	if startIdx < -1 {
		return nil, fmt.Errorf("invalid groups index %d", startIdx)
	}

	endIdx := int(totalGroupsNum - 1)
	if page > totalPagesNum {
		return nil, fmt.Errorf("offset page(%d) out of range", page)
	} else if totalPagesNum > page {
		// 计算结束索引
		endIdx = int(page*recordsPerPage) - 1
		// 因加入了all组进行计算，结束索引需往前移动一位
		endIdx--
	}

	for ; startIdx <= endIdx; startIdx++ {
		if startIdx == -1 {
			groups.GroupsMap["all"] = i.GenerateGroupAll()
			continue
		}
		groupName := i.sortedGroupNames[startIdx]
		groups.GroupsMap[groupName] = i.groups[groupName]
	}

	*groups.TotalGroupsNum = totalGroupsNum
	*groups.TotalPagesNum = totalPagesNum
	return groups, nil
}

func (i inventory) GetAllGroups() map[string]Group {
	return i.groups
}

func (i *inventory) GetTruncatedGroup() map[string]bool {
	truncatedGroup := make(map[string]bool)
	for groupName := range i.isTruncatedGroup {
		if i.isTruncatedGroup[groupName] {
			truncatedGroup[groupName] = true
		}
	}
	return truncatedGroup
}

func (i inventory) GenerateGroupAll() Group {
	groupAll := newGroup()
	_ = groupAll.setName("all")
	count := 0
	for _, g := range i.groups {
		_ = groupAll.AddHost(g.GetHosts()...)
		count += g.HostsLen()
		// todo: handle error
		//err := groupAll.AddHost(g.GetHosts()...)
		//if err != nil {
		//	fmt.Printf("get hosts from group %s failed, cased by: %s\n", g.GetName(), err)
		//}
	}
	return groupAll
}

func (i *inventory) sortGroups() {
	sort.Slice(i.sortedGroupNames, func(x, y int) bool {
		return isLessString(i.sortedGroupNames[x], i.sortedGroupNames[y])
	})
}

func (i *inventory) searchGroup(groupName string) int {
	return sort.Search(len(i.sortedGroupNames), func(idx int) bool {
		// ! idxName < groupName
		return !isLessString(i.sortedGroupNames[idx], groupName)
	})
}

func isLessString(x, y string) bool {
	// DONE: 数字字符串按数值大小排序
	x, y = strings.ToLower(x), strings.ToLower(y)
	xLen := len(x)
	yLen := len(y)
	var minLen int
	if xLen < yLen {
		minLen = xLen
	} else {
		minLen = yLen
	}

	for j := 0; j < minLen; j++ {
		switch {
		case x[j] == y[j]:
			continue
		case isNum(x[j]):
			switch {
			case isNum(y[j]):
				return isLessNumHead(x[j:], y[j:])
			default:
				return true
			}
		case isNum(y[j]):
			switch {
			case isNum(x[j]):
				return isLessNumHead(x[j:], y[j:])
			default:
				return false
			}
		default:
			return x[j] < y[j]
		}
	}
	return xLen == minLen && xLen != yLen
}

func isLessNumHead(x, y string) bool {
	n := len(x)
	m := len(y)
	var i, j int
	var xNumHead, yNumHead int
	for i < n || j < m {
		if i < n {
			if isNum(x[i]) {
				xNumHead *= 10
				xNumHead += int(x[i] - '0')
				i++
			} else {
				i = n
			}
		}
		if j < m {
			if isNum(y[j]) {
				yNumHead *= 10
				yNumHead += int(y[j] - '0')
				j++
			} else {
				j = m
			}
		}
	}
	return xNumHead < yNumHead
}

func isNum(s byte) bool {
	return s >= '0' && s <= '9'
}
