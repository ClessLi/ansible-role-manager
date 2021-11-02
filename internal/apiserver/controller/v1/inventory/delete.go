package inventory

import (
	ctrl_v1 "github.com/ClessLi/ansible-role-manager/internal/apiserver/controller/v1"
	metav1 "github.com/ClessLi/ansible-role-manager/internal/pkg/meta/v1"
	"github.com/gin-gonic/gin"
)

func (i *InventoryController) Delete(c *gin.Context) {
	if err := i.srv.Inventory().Delete(c, c.Param("name"), metav1.DeleteOptions{}); err != nil {
		ctrl_v1.WriteResponse(c, err, nil)

		return
	}

	ctrl_v1.WriteResponse(c, nil, nil)
}
