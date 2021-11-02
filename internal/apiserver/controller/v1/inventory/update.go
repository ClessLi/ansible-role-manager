package inventory

import (
	v1 "github.com/ClessLi/ansible-role-manager/api/apiserver/v1"
	ctrl_v1 "github.com/ClessLi/ansible-role-manager/internal/apiserver/controller/v1"
	"github.com/ClessLi/ansible-role-manager/internal/pkg/code"
	metav1 "github.com/ClessLi/ansible-role-manager/internal/pkg/meta/v1"
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"
)

func (i *InventoryController) Update(c *gin.Context) {
	var r v1.Group
	if err := c.ShouldBindJSON(&r); err != nil {
		ctrl_v1.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)

		return
	}

	groupBO, err := decoderIns.DecodeGroup(&r)
	if err != nil {
		ctrl_v1.WriteResponse(c, errors.WithCode(code.ErrDecodingFailed, err.Error()), nil)

		return
	}

	if err = i.srv.Inventory().Update(c, groupBO, metav1.UpdateOptions{}); err != nil {
		ctrl_v1.WriteResponse(c, err, nil)

		return
	}

	ctrl_v1.WriteResponse(c, nil, r)
}
