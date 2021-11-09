package inventory

import (
	ctrl_v1 "github.com/ClessLi/ansible-role-manager/internal/apiserver/controller/v1"
	metav1 "github.com/ClessLi/ansible-role-manager/internal/pkg/meta/v1"
	"github.com/gin-gonic/gin"
)

func (i *InventoryController) Get(c *gin.Context) {
	if groupBO, err := i.srv.Inventory().Get(c, c.Param("name"), metav1.GetOptions{}); err != nil {
		ctrl_v1.WriteResponse(c, err, nil)

		return
	} else {
		rsp := encoderIns.EncodeGroup(groupBO)

		ctrl_v1.WriteResponse(c, nil, rsp)
	}
}
