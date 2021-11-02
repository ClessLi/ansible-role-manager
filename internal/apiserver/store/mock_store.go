package store

import (
	"context"
	ansible_inventory "github.com/ClessLi/ansible-role-manager/internal/pkg/ansible-inventory"
	metav1 "github.com/ClessLi/ansible-role-manager/internal/pkg/meta/v1"
	"github.com/golang/mock/gomock"
	"reflect"
)

type MockFactory struct {
	ctrl     *gomock.Controller
	recorder *MockFactoryMockRecorder
}

type MockFactoryMockRecorder struct {
	moke *MockFactory
}

func NewMockFactory(ctrl *gomock.Controller) *MockFactory {
	mock := &MockFactory{ctrl: ctrl}
	mock.recorder = &MockFactoryMockRecorder{moke: mock}
	return mock
}

func (m *MockFactory) EXPECT() *MockFactoryMockRecorder {
	return m.recorder
}

func (m *MockFactory) Inventory() InventoryStore {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Inventory")
	ret0, _ := ret[0].(InventoryStore)
	return ret0
}

func (mr *MockFactoryMockRecorder) Inventory() *gomock.Call {
	mr.moke.ctrl.T.Helper()
	return mr.moke.ctrl.RecordCallWithMethodType(mr.moke, "Inventory", reflect.TypeOf((*MockFactory)(nil).Inventory))
}

type MockInventoryStore struct {
	ctrl     *gomock.Controller
	recorder *MockInventoryStoreMockRecorder
}

type MockInventoryStoreMockRecorder struct {
	mock *MockInventoryStore
}

func NewMockInventoryStore(ctrl *gomock.Controller) *MockInventoryStore {
	mock := &MockInventoryStore{ctrl: ctrl}
	mock.recorder = &MockInventoryStoreMockRecorder{mock: mock}
	return mock
}

func (m *MockInventoryStore) EXPECT() *MockInventoryStoreMockRecorder {
	return m.recorder
}

func (m *MockInventoryStore) Create(arg0 context.Context, arg1 ansible_inventory.Group, arg2 metav1.CreateOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockInventoryStoreMockRecorder) Create(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockInventoryStore)(nil).Create), arg0, arg1, arg2)
}

func (m *MockInventoryStore) Delete(arg0 context.Context, arg1 string, arg2 metav1.DeleteOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockInventoryStoreMockRecorder) Delete(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockInventoryStore)(nil).Delete), arg0, arg1, arg2)
}

func (m *MockInventoryStore) DeleteCollection(arg0 context.Context, arg1 []string, arg2 metav1.DeleteOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCollection", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockInventoryStoreMockRecorder) DeleteCollection(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCollection", reflect.TypeOf((*MockInventoryStore)(nil).DeleteCollection), arg0, arg1, arg2)
}

func (m *MockInventoryStore) Update(arg0 context.Context, arg1 ansible_inventory.Group, arg2 metav1.UpdateOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockInventoryStoreMockRecorder) Update(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockInventoryStore)(nil).Update), arg0, arg1, arg2)
}

func (m *MockInventoryStore) Get(arg0 context.Context, arg1 string, arg2 metav1.GetOptions) (ansible_inventory.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1, arg2)
	ret0, _ := ret[0].(ansible_inventory.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockInventoryStoreMockRecorder) Get(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockInventoryStore)(nil).Get), arg0, arg1, arg2)
}

func (m *MockInventoryStore) List(arg0 context.Context, arg1 metav1.ListOptions) (*ansible_inventory.Groups, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].(*ansible_inventory.Groups)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockInventoryStoreMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockInventoryStore)(nil).List), arg0, arg1)
}
