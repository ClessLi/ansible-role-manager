package inventory

import (
	"github.com/AlekSi/pointer"
	v1 "github.com/ClessLi/ansible-role-manager/api/apiserver/v1"
	ansible_inventory "github.com/ClessLi/ansible-role-manager/internal/pkg/ansible-inventory"
	"github.com/ClessLi/ansible-role-manager/internal/pkg/code"
	"github.com/marmotedu/errors"
)

var decoderIns = NewDecoder()

type Decoder interface {
	DecodeGroup(groupVO *v1.Group) (ansible_inventory.Group, error)
	DecodeGroups(groupsVO *v1.Groups) (*ansible_inventory.Groups, error)
}

type decoder struct {
}

func NewDecoder() Decoder {
	return Decoder(new(decoder))
}

func (d *decoder) DecodeGroup(groupVO *v1.Group) (ansible_inventory.Group, error) {
	hostsBO := make([]ansible_inventory.Host, 0)
	invalidHosts := make([]string, 0)
	for _, host := range groupVO.Hosts {
		hostBO := ansible_inventory.ParseHost(host.Ipaddr)
		if hostBO == nil {
			invalidHosts = append(invalidHosts, host.Ipaddr)
			continue
		}
		hostsBO = append(hostsBO, hostBO)
	}

	if len(invalidHosts) > 0 {
		return nil, errors.WithCode(code.ErrDecodingFailed, "invalid hosts: %v", invalidHosts)
	}

	groupBO, err := ansible_inventory.NewGroup(groupVO.GroupName, hostsBO)
	if err != nil {
		return nil, errors.WithCode(code.ErrDecodingFailed, err.Error())
	}
	return groupBO, nil
}

func (d *decoder) DecodeGroups(groupsVO *v1.Groups) (*ansible_inventory.Groups, error) {
	groupsBO := &ansible_inventory.Groups{GroupsMap: make(map[string]ansible_inventory.Group)}
	groupsBO.TotalGroupsNum = pointer.ToUint(groupsVO.TotalGroupsNum)
	groupsBO.TotalPagesNum = pointer.ToUint(groupsVO.TotalPagesNum)

	for groupName, groupVO := range groupsVO.Items {
		groupBO, err := d.DecodeGroup(groupVO)
		if err != nil {
			return groupsBO, err
		}
		groupsBO.GroupsMap[groupName] = groupBO
	}

	return groupsBO, nil
}
