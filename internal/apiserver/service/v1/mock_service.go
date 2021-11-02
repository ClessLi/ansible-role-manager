package v1

import (
	"context"
	ansible_inventory "github.com/ClessLi/ansible-role-manager/internal/pkg/ansible-inventory"
	metav1 "github.com/ClessLi/ansible-role-manager/internal/pkg/meta/v1"
	"github.com/golang/mock/gomock"
	"reflect"
)

type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

type MockServiceMockRecorder struct {
	mock *MockService
}

func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock: mock}
	return mock
}

func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

func (m *MockService) Inventory() InventorySrv {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Inventory")
	ret0, _ := ret[0].(InventorySrv)
	return ret0
}

func (mr *MockServiceMockRecorder) Inventory() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Inventory", reflect.TypeOf((*MockService)(nil).Inventory))
}

type MockInventorySrv struct {
	ctrl     *gomock.Controller
	recorder *MockInventorySrvMockRecorder
}

type MockInventorySrvMockRecorder struct {
	mock *MockInventorySrv
}

func NewMockInventorySrv(ctrl *gomock.Controller) *MockInventorySrv {
	mock := &MockInventorySrv{ctrl: ctrl}
	mock.recorder = &MockInventorySrvMockRecorder{mock: mock}
	return mock
}

func (m *MockInventorySrv) EXPECT() *MockInventorySrvMockRecorder {
	return m.recorder
}

func (m *MockInventorySrv) Create(arg0 context.Context, arg1 ansible_inventory.Group, arg2 metav1.CreateOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockInventorySrvMockRecorder) Create(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockInventorySrv)(nil).Create), arg0, arg1, arg2)
}

func (m *MockInventorySrv) Delete(arg0 context.Context, arg1 string, arg2 metav1.DeleteOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockInventorySrvMockRecorder) Delete(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockInventorySrv)(nil).Delete), arg0, arg1, arg2)
}

func (m *MockInventorySrv) DeleteCollection(arg0 context.Context, arg1 []string, arg2 metav1.DeleteOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCollection", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockInventorySrvMockRecorder) DeleteCollection(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCollection", reflect.TypeOf((*MockInventorySrv)(nil).DeleteCollection), arg0, arg1, arg2)
}

func (m *MockInventorySrv) Update(arg0 context.Context, arg1 ansible_inventory.Group, arg2 metav1.UpdateOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockInventorySrvMockRecorder) Update(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockInventorySrv)(nil).Update), arg0, arg1, arg2)
}

func (m *MockInventorySrv) Get(arg0 context.Context, arg1 string, arg2 metav1.GetOptions) (ansible_inventory.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1, arg2)
	ret0, _ := ret[0].(ansible_inventory.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockInventorySrvMockRecorder) Get(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockInventorySrv)(nil).Get), arg0, arg1, arg2)
}

func (m *MockInventorySrv) List(arg0 context.Context, arg1 metav1.ListOptions) (*ansible_inventory.Groups, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].(*ansible_inventory.Groups)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockInventorySrvMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockInventorySrv)(nil).List), arg0, arg1)
}
