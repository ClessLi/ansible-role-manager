package inventory

import (
	ctrl_v1 "github.com/ClessLi/ansible-role-manager/internal/apiserver/controller/v1"
	metav1 "github.com/ClessLi/ansible-role-manager/internal/pkg/meta/v1"
	"github.com/gin-gonic/gin"
)

func (i *InventoryController) DeleteCollection(c *gin.Context) {
	groupNames := c.QueryArray("name")

	if err := i.srv.Inventory().DeleteCollection(c, groupNames, metav1.DeleteOptions{}); err != nil {
		ctrl_v1.WriteResponse(c, err, nil)

		return
	}

	ctrl_v1.WriteResponse(c, nil, nil)
}
