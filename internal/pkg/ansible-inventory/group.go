package ansible_inventory

import (
	"errors"
	"sort"
	"strings"
)

type Groups struct {
	TotalGroupsNum *uint            `json:"totalGroupsNum"`
	TotalPagesNum  *uint            `json:"totalPagesNum"`
	GroupsMap      map[string]Group `json:"groupsMap"`
}

type Group interface {
	AddHost(hosts ...Host) error
	RemoveHost(hosts ...Host)
	setName(newName string) error
	GetName() string

	// GetHosts TODOList:
	// DONE: 2.hostsPattern 结构体实现，提供主机ip段存放，及提供数量统计
	GetHosts() []Host
	// HostsLen TODOList:
	// DONE: 1.hosts 数量统计
	HostsLen() int
}

type group struct {
	groupName string
	hosts     []Host
}

func NewGroup(groupName string, hosts []Host) (Group, error) {
	g := newGroup()
	err := g.setName(groupName)
	if err != nil {
		return nil, err
	}
	err = g.AddHost(hosts...)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func newGroup() Group {
	g := &group{hosts: make([]Host, 0)}
	return Group(g)
}

func (g *group) AddHost(hosts ...Host) error {
	if hosts == nil || len(hosts) == 0 {
		return errors.New("input hosts are null")
	}
	isIn := true
	for _, h := range hosts {
		if h == nil {
			continue
		}
		idx := sort.Search(len(g.hosts), func(i int) bool {
			return !g.hosts[i].Less(h) || g.hosts[i].IsInclude(h) || h.IsInclude(g.hosts[i])
		})
		if idx < len(g.hosts) {
			// 判断索引对象与插入对象是否相同或者涵盖插入对象
			if g.hosts[idx].Equal(h) || g.hosts[idx].IsInclude(h) {
				continue
			}
			// 判断插入对象是否涵盖索引对象
			if h.IsInclude(g.hosts[idx]) {
				// 如果涵盖则删除索引对象
				g.hosts = append(g.hosts[:idx], g.hosts[idx+1:]...)
			}
		}
		isIn = false
		g.hosts = append(g.hosts, h)
		sort.Slice(g.hosts, func(x, y int) bool {
			return g.hosts[x].Less(g.hosts[y])
		})
	}
	if isIn {
		return errors.New("input hosts are in the group")
	}
	return nil
}

func (g *group) RemoveHost(hosts ...Host) {
	if hosts == nil || len(hosts) == 0 {
		return
	}
	for _, h := range hosts {
		if h == nil {
			continue
		}
		idx := sort.Search(len(g.hosts), func(i int) bool {
			return g.hosts[i].Less(h)
		})
		if idx < len(g.hosts) && g.hosts[idx].Equal(h) {
			g.hosts = append(g.hosts[:idx], g.hosts[idx+1:]...)
		}
	}
}

func (g *group) setName(newName string) error {
	newName = strings.TrimSpace(newName)
	if newName == "" {
		return errors.New("can not renew group name to null")
	}
	g.groupName = newName
	return nil
}

func (g group) GetName() string {
	return g.groupName
}

func (g group) GetHosts() []Host {
	return g.hosts
}

func (g group) HostsLen() int {
	num := 0
	for _, h := range g.hosts {
		num += h.Len()
	}
	return num
}
