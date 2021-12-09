package inventory

import (
	"github.com/ClessLi/ansible-role-manager/internal/pkg/core"
	metav1 "github.com/ClessLi/ansible-role-manager/internal/pkg/meta/v1"
	log "github.com/ClessLi/ansible-role-manager/pkg/log/v2"
	"github.com/gin-gonic/gin"
)

func (i *InventoryController) DeleteCollection(c *gin.Context) {
	log.L(c).Info("batch delete inventory group function called.")

	groupNames := c.QueryArray("groups")

	if err := i.srv.Inventory().DeleteCollection(c, groupNames, metav1.DeleteOptions{}); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
