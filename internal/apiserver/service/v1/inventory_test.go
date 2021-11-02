package v1

import (
	"context"
	"fmt"
	"github.com/ClessLi/ansible-role-manager/internal/apiserver/store"
	"github.com/ClessLi/ansible-role-manager/internal/apiserver/store/fake"
	ansible_inventory "github.com/ClessLi/ansible-role-manager/internal/pkg/ansible-inventory"
	v1 "github.com/ClessLi/ansible-role-manager/internal/pkg/meta/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

type Suite struct {
	suite.Suite
	mockFactory *store.MockFactory

	mockInventoryStore *store.MockInventoryStore
	groupsMap          map[string]ansible_inventory.Group
}

func (s *Suite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.groupsMap = fake.FakeGroups(10)
	s.mockFactory = store.NewMockFactory(ctrl)
	s.mockInventoryStore = store.NewMockInventoryStore(ctrl)
	s.mockFactory.EXPECT().Inventory().AnyTimes().Return(s.mockInventoryStore)
}

func TestInventory(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test_inventoryService_Create() {
	canCreateGroupName := "test-group1"
	existGroupName := "test-group3"
	// can Create group: test-group1
	s.mockInventoryStore.EXPECT().Create(gomock.Any(), gomock.Eq(s.groupsMap[canCreateGroupName]), gomock.Any()).Return(nil)
	// can not Create the exist group: test-group3
	s.mockInventoryStore.EXPECT().Create(gomock.Any(), gomock.Eq(s.groupsMap[existGroupName]), gomock.Any()).Return(fmt.Errorf("group '%v' is exist", existGroupName))

	type fields struct {
		store store.Factory
	}
	type args struct {
		ctx     context.Context
		group   ansible_inventory.Group
		options v1.CreateOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "normal test",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				ctx:     context.TODO(),
				group:   s.groupsMap[canCreateGroupName],
				options: v1.CreateOptions{},
			},
			wantErr: false,
		},
		{
			name: "create with exist group",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				ctx:     context.TODO(),
				group:   s.groupsMap[existGroupName],
				options: v1.CreateOptions{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			i := &inventoryService{
				store: tt.fields.store,
			}
			if err := i.Create(tt.args.ctx, tt.args.group, tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (s *Suite) Test_inventoryService_Delete() {
	canDeleteGroupName := "test-group3"
	nonexistentGroupName := "test-group1"
	// can Delete group: test-group3
	s.mockInventoryStore.EXPECT().Delete(gomock.Any(), gomock.Eq(canDeleteGroupName), gomock.Any()).Return(nil)
	// can not Delete the nonexistent group: test-group1
	s.mockInventoryStore.EXPECT().Delete(gomock.Any(), gomock.Eq(nonexistentGroupName), gomock.Any()).Return(fmt.Errorf("group '%v' is not exist", nonexistentGroupName))
	type fields struct {
		store store.Factory
	}
	type args struct {
		ctx       context.Context
		groupName string
		options   v1.DeleteOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "normal test",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				ctx:       context.TODO(),
				groupName: canDeleteGroupName,
				options:   v1.DeleteOptions{},
			},
			wantErr: false,
		},
		{
			name: "delete with nonexistent group",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				ctx:       context.TODO(),
				groupName: nonexistentGroupName,
				options:   v1.DeleteOptions{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			i := &inventoryService{
				store: tt.fields.store,
			}
			if err := i.Delete(tt.args.ctx, tt.args.groupName, tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (s *Suite) Test_inventoryService_DeleteCollection() {
	canDeleteGroup2Name := "test-group2"
	canDeleteGroup3Name := "test-group3"
	nonexistentGroupName := "test-group1"
	// can DeleteCollection group in [ test-group2, test-group3 ]
	s.mockInventoryStore.EXPECT().DeleteCollection(gomock.Any(), gomock.Eq([]string{canDeleteGroup2Name, canDeleteGroup3Name}), gomock.Any()).Return(nil)
	// can not DeleteCollection is one of test-group1
	s.mockInventoryStore.EXPECT().DeleteCollection(gomock.Any(),
		gomock.All(
			gomock.Not([]string{canDeleteGroup2Name}),
			gomock.Not([]string{canDeleteGroup3Name}),
			gomock.Not([]string{canDeleteGroup2Name, canDeleteGroup3Name}),
		),
		gomock.Any()).Return(fmt.Errorf("nonexistent group list: %+v", []string{nonexistentGroupName}))

	type fields struct {
		store store.Factory
	}
	type args struct {
		ctx        context.Context
		groupNames []string
		options    v1.DeleteOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "normal delete collection",
			fields: fields{store: s.mockFactory},
			args: args{
				ctx:        context.TODO(),
				groupNames: []string{canDeleteGroup2Name, canDeleteGroup3Name},
				options:    v1.DeleteOptions{},
			},
			wantErr: false,
		},
		{
			name:   "delete collection with any nonexistent groups list",
			fields: fields{store: s.mockFactory},
			args: args{
				ctx:        context.TODO(),
				groupNames: []string{canDeleteGroup2Name, nonexistentGroupName},
				options:    v1.DeleteOptions{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			i := &inventoryService{
				store: tt.fields.store,
			}
			if err := i.DeleteCollection(tt.args.ctx, tt.args.groupNames, tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("DeleteCollection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_inventoryService_Get(t *testing.T) {
	type fields struct {
		store store.Factory
	}
	type args struct {
		ctx       context.Context
		groupName string
		options   v1.GetOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ansible_inventory.Group
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &inventoryService{
				store: tt.fields.store,
			}
			got, err := i.Get(tt.args.ctx, tt.args.groupName, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inventoryService_List(t *testing.T) {
	type fields struct {
		store store.Factory
	}
	type args struct {
		ctx     context.Context
		options v1.ListOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ansible_inventory.Groups
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &inventoryService{
				store: tt.fields.store,
			}
			got, err := i.List(tt.args.ctx, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inventoryService_Update(t *testing.T) {
	type fields struct {
		store store.Factory
	}
	type args struct {
		ctx     context.Context
		group   ansible_inventory.Group
		options v1.UpdateOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &inventoryService{
				store: tt.fields.store,
			}
			if err := i.Update(tt.args.ctx, tt.args.group, tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
