package file

import (
	"context"
	"fmt"
	"github.com/AlekSi/pointer"
	"github.com/ClessLi/ansible-role-manager/internal/apiserver/store"
	ansible_inventory "github.com/ClessLi/ansible-role-manager/internal/pkg/ansible-inventory"
	file_store "github.com/ClessLi/ansible-role-manager/internal/pkg/file-store"
	v1 "github.com/ClessLi/ansible-role-manager/internal/pkg/meta/v1"
	"path/filepath"
	"reflect"
	"testing"
)

type mockFileStore struct {
}

func (m mockFileStore) Write(filePath string, data []byte) error {
	return nil
}

func (m mockFileStore) Read(filePath string) (data []byte, err error) {
	groups := testGroupExample()
	_, groupName := filepath.Split(filePath)
	if group, has := groups[groupName]; has {
		return ansible_inventory.NewParser().Dump(group)
	}
	return nil, fmt.Errorf("mock FileStore: %v is not exist", filePath)
}

func (m mockFileStore) Remove(filePath string) error {
	return nil
}

func (m mockFileStore) AllFiles() ([]string, error) {
	fileList := make([]string, 0)
	groups := testGroupExample()
	for groupName := range groups {
		fileList = append(fileList, groupName)
	}
	return fileList, nil
}

func (m mockFileStore) Workspace() (string, error) {
	return "mockPath", nil
}

func testGroupExample() map[string]ansible_inventory.Group {
	group1, _ := ansible_inventory.NewGroup("test-group", []ansible_inventory.Host{ansible_inventory.ParseHost("192.168.0.1"), ansible_inventory.ParseHost("10.1.0.1"), ansible_inventory.ParseHost("192.168.[11:100].[1:254]")})
	group2, _ := ansible_inventory.NewGroup("test-group2", []ansible_inventory.Host{ansible_inventory.ParseHost("192.168.1.1"), ansible_inventory.ParseHost("10.2.0.1"), ansible_inventory.ParseHost("192.168.[21:200].[1:254]")})
	groups := make(map[string]ansible_inventory.Group)
	groups[group1.GetName()] = group1
	groups[group2.GetName()] = group2
	return groups
}

func testNewGroup3() ansible_inventory.Group {
	group, _ := ansible_inventory.NewGroup("test-group3", []ansible_inventory.Host{ansible_inventory.ParseHost("192.168.2.1"), ansible_inventory.ParseHost("10.3.0.1"), ansible_inventory.ParseHost("192.168.[31:150].[1:254]")})
	return group
}

func Test_inventory_load(t *testing.T) {
	dir := "../../../../test/testdata/inventory"
	fs, err := file_store.NewFileStore(dir)
	if err != nil {
		panic(err)
	}
	nullFilesfs, err := file_store.NewFileStore(filepath.Dir(dir))
	if err != nil {
		panic(err)
	}
	groups := testGroupExample()
	parser := ansible_inventory.NewParser()
	type fields struct {
		fs     file_store.FileStore
		parser ansible_inventory.Parser
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ansible_inventory.Inventory
		wantErr bool
	}{
		{
			name: "normal load",
			fields: fields{
				fs:     fs,
				parser: parser,
			},
			args: args{ctx: context.Background()},
			want: ansible_inventory.NewInventory(groups),
		},
		//{
		//	name: "error dir path",
		//	fields: fields{
		//		fs:    "testdir",
		//		parser: ansible_inventory.NewParser(),
		//	},
		//	wantErr: true,
		//},
		{
			name: "null inventory file dir",
			fields: fields{
				fs:     nullFilesfs,
				parser: parser,
			},
			args:    args{ctx: context.Background()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &inventory{
				fs:     tt.fields.fs,
				parser: tt.fields.parser,
			}
			got, err := i.load(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("load() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inventory_save(t *testing.T) {
	dir := "../../../../test/testdata/inventory"
	fs, err := file_store.NewFileStore(dir)
	if err != nil {
		panic(err)
	}
	groups := testGroupExample()
	type fields struct {
		fs     file_store.FileStore
		parser ansible_inventory.Parser
	}
	type args struct {
		ctx context.Context
		inv ansible_inventory.Inventory
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "normal save",
			fields: fields{
				fs:     fs,
				parser: ansible_inventory.NewParser(),
			},
			args: args{
				ctx: context.Background(),
				inv: ansible_inventory.NewInventory(groups),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &inventory{
				fs:     tt.fields.fs,
				parser: tt.fields.parser,
			}
			if err := i.save(tt.args.ctx, tt.args.inv); (err != nil) != tt.wantErr {
				t.Errorf("save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetFileFactory(t *testing.T) {
	tests := []struct {
		name    string
		want    store.Factory
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFileFactory()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileFactory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFileFactory() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileStore_Inventory(t *testing.T) {
	type fields struct {
		fs file_store.FileStore
	}
	tests := []struct {
		name   string
		fields fields
		want   store.InventoryStore
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fileStore{
				fs: tt.fields.fs,
			}
			if got := f.Inventory(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Inventory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inventory_Create(t *testing.T) {

	type fields struct {
		fs     file_store.FileStore
		parser ansible_inventory.Parser
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
				fs:     new(mockFileStore),
				parser: ansible_inventory.NewParser(),
			},
			args: args{
				ctx:     context.Background(),
				group:   testNewGroup3(),
				options: v1.CreateOptions{},
			},
		},
		{
			name: "create exist group",
			fields: fields{
				fs:     new(mockFileStore),
				parser: ansible_inventory.NewParser(),
			},
			args: args{
				ctx:     context.Background(),
				group:   testGroupExample()["test-group"],
				options: v1.CreateOptions{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &inventory{
				fs:     tt.fields.fs,
				parser: tt.fields.parser,
			}
			if err := i.Create(tt.args.ctx, tt.args.group, tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_inventory_Delete(t *testing.T) {
	type fields struct {
		fs     file_store.FileStore
		parser ansible_inventory.Parser
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
				fs:     new(mockFileStore),
				parser: ansible_inventory.NewParser(),
			},
			args: args{
				ctx:       context.Background(),
				groupName: "test-group",
				options:   v1.DeleteOptions{Force: false},
			},
		},
		{
			name: "delete nonexistent group",
			fields: fields{
				fs:     new(mockFileStore),
				parser: ansible_inventory.NewParser(),
			},
			args: args{
				ctx:       context.Background(),
				groupName: "nonexistent-group",
				options:   v1.DeleteOptions{Force: false},
			},
			wantErr: true,
		},
		{
			name: "delete nonexistent group with Force",
			fields: fields{
				fs:     new(mockFileStore),
				parser: ansible_inventory.NewParser(),
			},
			args: args{
				ctx:       context.Background(),
				groupName: "nonexistent-group",
				options:   v1.DeleteOptions{Force: true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &inventory{
				fs:     tt.fields.fs,
				parser: tt.fields.parser,
			}
			if err := i.Delete(tt.args.ctx, tt.args.groupName, tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_inventory_DeleteCollection(t *testing.T) {
	type fields struct {
		fs     file_store.FileStore
		parser ansible_inventory.Parser
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
			name: "normal test",
			fields: fields{
				fs:     new(mockFileStore),
				parser: ansible_inventory.NewParser(),
			},
			args: args{
				ctx:        context.Background(),
				groupNames: []string{"test-group", "test-group2"},
				options:    v1.DeleteOptions{Force: false},
			},
		},
		{
			name: "delete groups, which has some nonexistent groups",
			fields: fields{
				fs:     new(mockFileStore),
				parser: ansible_inventory.NewParser(),
			},
			args: args{
				ctx:        context.Background(),
				groupNames: []string{"test-group", "test-group3", "nonexistent-group"},
				options:    v1.DeleteOptions{Force: false},
			},
			wantErr: true,
		},
		{
			name: "delete groups with Force, which has some nonexistent groups",
			fields: fields{
				fs:     new(mockFileStore),
				parser: ansible_inventory.NewParser(),
			},
			args: args{
				ctx:        context.Background(),
				groupNames: []string{"test-group", "test-group3", "nonexistent-group"},
				options:    v1.DeleteOptions{Force: true},
			},
		},
		{
			name: "delete nil groups",
			fields: fields{
				fs:     new(mockFileStore),
				parser: ansible_inventory.NewParser(),
			},
			args: args{
				ctx:        context.Background(),
				groupNames: nil,
				options:    v1.DeleteOptions{Force: false},
			},
			wantErr: true,
		},
		{
			name: "delete null groups",
			fields: fields{
				fs:     new(mockFileStore),
				parser: ansible_inventory.NewParser(),
			},
			args: args{
				ctx:        context.Background(),
				groupNames: []string{},
				options:    v1.DeleteOptions{Force: false},
			},
			wantErr: true,
		},
		{
			name: "delete nil groups with Force",
			fields: fields{
				fs:     new(mockFileStore),
				parser: ansible_inventory.NewParser(),
			},
			args: args{
				ctx:        context.Background(),
				groupNames: nil,
				options:    v1.DeleteOptions{Force: true},
			},
		},
		{
			name: "delete null groups with Force",
			fields: fields{
				fs:     new(mockFileStore),
				parser: ansible_inventory.NewParser(),
			},
			args: args{
				ctx:        context.Background(),
				groupNames: []string{},
				options:    v1.DeleteOptions{Force: true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &inventory{
				fs:     tt.fields.fs,
				parser: tt.fields.parser,
			}
			if err := i.DeleteCollection(tt.args.ctx, tt.args.groupNames, tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("DeleteCollection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_inventory_Update(t *testing.T) {
	deltaGroup2 := testGroupExample()["test-group2"]
	deltaGroup2.RemoveHost(deltaGroup2.GetHosts()[0])
	err := deltaGroup2.AddHost(ansible_inventory.ParseHost("10.10.10.[1:255]"))
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		fs     file_store.FileStore
		parser ansible_inventory.Parser
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
		{
			name: "normal test",
			fields: fields{
				fs:     new(mockFileStore),
				parser: ansible_inventory.NewParser(),
			},
			args: args{
				ctx:     context.Background(),
				group:   deltaGroup2,
				options: v1.UpdateOptions{},
			},
		},
		{
			name: "update a nonexistent group",
			fields: fields{
				fs:     new(mockFileStore),
				parser: ansible_inventory.NewParser(),
			},
			args: args{
				ctx:     context.Background(),
				group:   testNewGroup3(),
				options: v1.UpdateOptions{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &inventory{
				fs:     tt.fields.fs,
				parser: tt.fields.parser,
			}
			if err := i.Update(tt.args.ctx, tt.args.group, tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_inventory_Get(t *testing.T) {
	type fields struct {
		fs     file_store.FileStore
		parser ansible_inventory.Parser
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
		{
			name: "normal test",
			fields: fields{
				fs:     new(mockFileStore),
				parser: ansible_inventory.NewParser(),
			},
			args: args{
				ctx:       context.Background(),
				groupName: "test-group",
				options:   v1.GetOptions{},
			},
			want: testGroupExample()["test-group"],
		},
		{
			name: "get a nonexistent group",
			fields: fields{
				fs:     new(mockFileStore),
				parser: ansible_inventory.NewParser(),
			},
			args: args{
				ctx:       context.Background(),
				groupName: "nonexistent-group",
				options:   v1.GetOptions{},
			},
			wantErr: true,
		},
		{
			name: "get group with null name( the same as 'get a nonexistent group' test",
			fields: fields{
				fs:     new(mockFileStore),
				parser: ansible_inventory.NewParser(),
			},
			args: args{
				ctx:       context.Background(),
				groupName: "",
				options:   v1.GetOptions{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &inventory{
				fs:     tt.fields.fs,
				parser: tt.fields.parser,
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

func Test_inventory_List(t *testing.T) {
	testInv := ansible_inventory.NewInventory(testGroupExample())

	onePtenNumListOPTS := v1.ListOptions{
		Page:       pointer.ToUint(1),
		NumPerPage: pointer.ToUint(10),
	}

	onePtenNumGroupsRet := &ansible_inventory.Groups{
		GroupsMap:      testGroupExample(),
		TotalGroupsNum: pointer.ToUint(2),
		TotalPagesNum:  pointer.ToUint(1),
	}
	onePtenNumGroupsRet.GroupsMap["all"] = testInv.GenerateGroupAll()

	twoPoneNumListOPTS := v1.ListOptions{
		Page:       pointer.ToUint(2),
		NumPerPage: pointer.ToUint(1),
	}

	twoPoneNumGroupsRet := &ansible_inventory.Groups{
		GroupsMap:      map[string]ansible_inventory.Group{"test-group": testGroupExample()["test-group"]},
		TotalGroupsNum: pointer.ToUint(2),
		TotalPagesNum:  pointer.ToUint(3),
	}

	twoPtenNumListOPTS := v1.ListOptions{
		Page:       pointer.ToUint(2),
		NumPerPage: pointer.ToUint(10),
	}

	type fields struct {
		fs     file_store.FileStore
		parser ansible_inventory.Parser
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
		{
			name: "normal test",
			fields: fields{
				fs:     new(mockFileStore),
				parser: ansible_inventory.NewParser(),
			},
			args: args{
				ctx:     context.Background(),
				options: onePtenNumListOPTS,
			},
			want: onePtenNumGroupsRet,
		},
		{
			name: "page 2, 1/page",
			fields: fields{
				fs:     new(mockFileStore),
				parser: ansible_inventory.NewParser(),
			},
			args: args{
				ctx:     context.Background(),
				options: twoPoneNumListOPTS,
			},
			want: twoPoneNumGroupsRet,
		},
		{
			name: "specify page is out of range",
			fields: fields{
				fs:     new(mockFileStore),
				parser: ansible_inventory.NewParser(),
			},
			args: args{
				ctx:     context.Background(),
				options: twoPtenNumListOPTS,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &inventory{
				fs:     tt.fields.fs,
				parser: tt.fields.parser,
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
