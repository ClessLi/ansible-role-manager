package inventory

import (
	v1 "github.com/ClessLi/ansible-role-manager/api/apiserver/v1"
	ansible_inventory "github.com/ClessLi/ansible-role-manager/internal/pkg/ansible-inventory"
)

var encoderIns = NewEncoder()

type Encoder interface {
	EncodeGroup(groupBO ansible_inventory.Group) *v1.Group
	EncodeGroups(groupsBO *ansible_inventory.Groups) *v1.Groups
}

type encoder struct {
}

func NewEncoder() Encoder {
	return Encoder(new(encoder))
}

func (e *encoder) EncodeGroup(groupBO ansible_inventory.Group) *v1.Group {
	hostsVO := make([]*v1.Host, 0)
	for _, hostBO := range groupBO.GetHosts() {
		hostsVO = append(hostsVO, &v1.Host{Ipaddr: hostBO.GetIPString()})
	}

	return &v1.Group{
		GroupName: groupBO.GetName(),
		Hosts:     hostsVO,
	}
}

func (e *encoder) EncodeGroups(groupsBO *ansible_inventory.Groups) *v1.Groups {
	groupsVO := &v1.Groups{
		TotalGroupsNum: *groupsBO.TotalGroupsNum,
		TotalPagesNum:  *groupsBO.TotalPagesNum,
		Items:          make(map[string]*v1.Group),
	}
	for groupName, groupBO := range groupsBO.GroupsMap {
		groupsVO.Items[groupName] = e.EncodeGroup(groupBO)
	}
	return groupsVO
}
