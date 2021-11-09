package inventory

import (
	ctrl_v1 "github.com/ClessLi/ansible-role-manager/internal/apiserver/controller/v1"
	"github.com/ClessLi/ansible-role-manager/internal/pkg/code"
	metav1 "github.com/ClessLi/ansible-role-manager/internal/pkg/meta/v1"
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"
)

func (i *InventoryController) List(c *gin.Context) {
	var r metav1.ListOptions
	if err := c.ShouldBindJSON(&r); err != nil {
		ctrl_v1.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)

		return
	}

	groupsBO, err := i.srv.Inventory().List(c, r)
	if err != nil {
		ctrl_v1.WriteResponse(c, err, nil)

		return
	}

	rst := encoderIns.EncodeGroups(groupsBO)
	ctrl_v1.WriteResponse(c, nil, rst)
}
