package inventory

import (
	"github.com/ClessLi/ansible-role-manager/internal/pkg/code"
	"github.com/ClessLi/ansible-role-manager/internal/pkg/core"
	metav1 "github.com/ClessLi/ansible-role-manager/internal/pkg/meta/v1"
	log "github.com/ClessLi/ansible-role-manager/pkg/log/v2"
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"
)

func (i *InventoryController) List(c *gin.Context) {
	log.L(c).Info("list inventory function called.")

	var r metav1.ListOptions
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)

		return
	}

	groupsBO, err := i.srv.Inventory().List(c, r)
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	rst := encoderIns.EncodeGroups(groupsBO)
	core.WriteResponse(c, nil, rst)
}
