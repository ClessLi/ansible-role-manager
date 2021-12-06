package inventory

import (
	"github.com/ClessLi/ansible-role-manager/internal/pkg/core"
	metav1 "github.com/ClessLi/ansible-role-manager/internal/pkg/meta/v1"
	"github.com/gin-gonic/gin"
)

func (i *InventoryController) Get(c *gin.Context) {
	if groupBO, err := i.srv.Inventory().Get(c, c.Param("group"), metav1.GetOptions{}); err != nil {
		core.WriteResponse(c, err, nil)

		return
	} else {
		rsp := encoderIns.EncodeGroup(groupBO)

		core.WriteResponse(c, nil, rsp)
	}
}
