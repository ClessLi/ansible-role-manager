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

func (i *InventoryController) Create(c *gin.Context) {
	log.L(c).Info("inventory group create function called.")

	var r v1.Group
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)

		return
	}

	groupBO, err := decoderIns.DecodeGroup(&r)
	if err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrDecodingFailed, err.Error()), nil)

		return
	}

	if err := i.srv.Inventory().Create(c, groupBO, metav1.CreateOptions{}); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, r)
}
