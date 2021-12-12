package inventory

import (
	v1 "github.com/ClessLi/ansible-role-manager/api/apiserver/v1"
	"github.com/ClessLi/ansible-role-manager/internal/pkg/code"
	"github.com/ClessLi/ansible-role-manager/internal/pkg/core"
	metav1 "github.com/ClessLi/ansible-role-manager/internal/pkg/meta/v1"
	log "github.com/ClessLi/ansible-role-manager/pkg/log/v2"
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"
)

func (i *InventoryController) Update(c *gin.Context) {
	log.L(c).Info("update inventory group function called.")

	var r v1.Group
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)

		return
	}

	groupBO, err := i.srv.Inventory().Get(c, c.Param("group"), metav1.GetOptions{})
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	r.GroupName = groupBO.GetName()

	requestGroupBO, err := decoderIns.DecodeGroup(&r)
	if err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrDecodingFailed, err.Error()), nil)

		return
	}

	if err = i.srv.Inventory().Update(c, requestGroupBO, metav1.UpdateOptions{}); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, r)
}
