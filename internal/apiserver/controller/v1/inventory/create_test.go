package inventory

import (
	"bytes"
	"fmt"
	srvv1 "github.com/ClessLi/ansible-role-manager/internal/apiserver/service/v1"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestInventoryController_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := srvv1.NewMockService(ctrl)
	mockInventorySrv := srvv1.NewMockInventorySrv(ctrl)
	mockInventorySrv.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockService.EXPECT().Inventory().Return(mockInventorySrv)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	body := bytes.NewBufferString(
		`{"group_name": "test-group1", "hosts": [{"ipaddr": "192.168.0.1"}, {"ipaddr": "10.0.0.1"}, {"ipaddr": "192.168.0.[1:254]"}]}`,
	)
	c.Request, _ = http.NewRequest("POST", "/v1/groups", body)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("requestID", fmt.Sprintf("test-%d", time.Now().Unix()))
	c.Set("username", "testuser")

	type fields struct {
		srv srvv1.Service
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "default",
			fields: fields{srv: mockService},
			args:   args{c: c},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InventoryController{
				srv: tt.fields.srv,
			}
			i.Create(tt.args.c)
		})
	}
}
