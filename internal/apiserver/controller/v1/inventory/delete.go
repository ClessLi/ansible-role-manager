package inventory

import (
	"github.com/ClessLi/ansible-role-manager/internal/pkg/core"
	metav1 "github.com/ClessLi/ansible-role-manager/internal/pkg/meta/v1"
	log "github.com/ClessLi/ansible-role-manager/pkg/log/v2"
	"github.com/gin-gonic/gin"
)

func (i *InventoryController) Delete(c *gin.Context) {
	log.L(c).Info("delete inventory group function called.")

	if err := i.srv.Inventory().Delete(c, c.Param("group"), metav1.DeleteOptions{}); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
