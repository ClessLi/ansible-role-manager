package ansible_inventory

import (
	"fmt"
	"reflect"
	"testing"
)

func testGroupExample() map[string]Group {
	group1 := newGroup()
	group2 := newGroup()
	_ = group1.setName("test-group")
	_ = group2.setName("test-group2")
	_ = group1.AddHost(ParseHost("192.168.0.1"), ParseHost("10.1.0.1"))
	_ = group2.AddHost(ParseHost("192.168.1.1"), ParseHost("10.2.0.1"))
	_ = group1.AddHost(ParseHost("192.168.[11:100].[1:254]"))
	_ = group2.AddHost(ParseHost("192.168.[21:200].[1:254]"))
	groups := make(map[string]Group)
	groups[group1.GetName()] = group1
	groups[group2.GetName()] = group2
	return groups
}

func Test_inventory_AddHostToGroup(t *testing.T) {
	type fields struct {
		groups           map[string]Group
		isTruncatedGroup map[string]bool
	}
	type args struct {
		groupName string
		hosts     []Host
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		addFailed bool
		wantErr   bool
	}{
		{
			name: "add host to exist group",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				groupName: "test-group",
				hosts:     []Host{ParseHost("192.168.1.2")},
			},
		},
		{
			name: "add hosts to exist group",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				groupName: "test-group2",
				hosts:     []Host{ParseHost("192.168.2.2"), ParseHost("10.2.0.2")},
			},
		},
		{
			name: "add hosts pattern to exist group",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				groupName: "test-group2",
				hosts:     []Host{ParseHost("192.168.[19:30].[1:100]"), ParseHost("10.2.3.[1:254]")},
			},
		},
		{
			name: "add host to not exist group",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				groupName: "test-group3",
				hosts:     []Host{ParseHost("192.168.3.1")},
			},
		},
		{
			name: "add hosts to not exist group",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				groupName: "test-group4",
				hosts:     []Host{ParseHost("192.168.4.1"), ParseHost("10.4.0.1")},
			},
		},
		{
			name: "add nil host to exist group",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				groupName: "test-group",
				hosts:     []Host{nil},
			},
			addFailed: true,
			wantErr:   true,
		},
		{
			name: "add nil hosts to exist group",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				groupName: "test-group",
				hosts:     nil,
			},
			addFailed: true,
			wantErr:   true,
		},
		{
			name: "add null host to exist group",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				groupName: "test-group",
				hosts:     []Host{},
			},
			addFailed: true,
			wantErr:   true,
		},
		{
			name: "add nil host to not exist group",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				groupName: "test-group3",
				hosts:     []Host{nil},
			},
			addFailed: true,
			wantErr:   true,
		},
		{
			name: "add nil hosts to not exist group",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				groupName: "test-group3",
				hosts:     nil,
			},
			addFailed: true,
			wantErr:   true,
		},
		{
			name: "add null host to not exist group",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				groupName: "test-group3",
				hosts:     []Host{},
			},
			addFailed: true,
			wantErr:   true,
		},
		{
			name: "add host to the group which named 'all'",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				groupName: "all",
				hosts:     []Host{ParseHost("192.168.4.1")},
			},
			addFailed: true,
			wantErr:   true,
		},
		{
			name: "add hosts to the group which named 'All'",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				groupName: "All",
				hosts:     []Host{ParseHost("192.168.[0-200].[1-254]"), ParseHost("10.1.4.100")},
			},
			addFailed: true,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &inventory{
				groups:           tt.fields.groups,
				isTruncatedGroup: tt.fields.isTruncatedGroup,
			}
			var bLen, aLen int
			if _, has := i.groups[tt.args.groupName]; has {
				bLen = i.groups[tt.args.groupName].HostsLen()
			}
			if err := i.AddHostToGroup(tt.args.groupName, tt.args.hosts...); (err != nil) != tt.wantErr {
				t.Errorf("AddHostGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, has := i.groups[tt.args.groupName]; has {
				aLen = i.groups[tt.args.groupName].HostsLen()
			}
			if (bLen >= aLen) != tt.addFailed {
				t.Errorf("AddHostToGroup() add failed, before lenghth = %d, after lenghth = %d", bLen, aLen)
			}
		})
	}
}

func Test_inventory_RenewGroupName(t *testing.T) {
	type fields struct {
		groups           map[string]Group
		isTruncatedGroup map[string]bool
	}
	type args struct {
		oldName string
		newName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "normal rename",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				oldName: "test-group",
				newName: "test-group1",
			},
		},
		{
			name: "not exist group",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				oldName: "test-group1",
				newName: "test-group3",
			},
			wantErr: true,
		},
		{
			name: "rename to duplicate group",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				oldName: "test-group",
				newName: "test-group2",
			},
			wantErr: true,
		},
		{
			name: "rename to 'all' group",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				oldName: "test-group",
				newName: "all",
			},
			wantErr: true,
		},
		{
			name: "rename to 'All' group",
			fields: fields{
				groups:           testGroupExample(),
				isTruncatedGroup: make(map[string]bool),
			},
			args: args{
				oldName: "test-group2",
				newName: "All",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := make(map[string]Group)
			for s, g2 := range tt.fields.groups {
				g[s] = newGroup()
				err := g[s].setName(g2.GetName())
				if err != nil {
					t.Fatal(err)
				}
				_ = g[s].AddHost(g2.GetHosts()...)
			}

			//i := &inventory{
			//	groups:           g,
			//	isTruncatedGroup: make(map[string]bool),
			//}
			i := NewInventory(g)
			if err := i.RenewGroupName(tt.args.oldName, tt.args.newName); (err != nil) != tt.wantErr {
				t.Errorf("RenewGroupName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				t.Logf("old group name = %v, want new group name = %v, got %v", tt.args.oldName, tt.args.newName, i.GetAllGroups()[tt.args.newName].GetName())
			}
		})
	}
}

func Test_isLessString(t *testing.T) {
	type args struct {
		x string
		y string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "abc < acb == true",
			args: args{
				x: "abc",
				y: "acb",
			},
			want: true,
		},
		{
			name: "abc < abc == false",
			args: args{
				x: "abc",
				y: "abc",
			},
			want: false,
		},
		{
			name: "abc < abca == true",
			args: args{
				x: "abc",
				y: "abca",
			},
			want: true,
		},
		{
			name: "abca < abc == false",
			args: args{
				x: "abca",
				y: "abc",
			},
			want: false,
		},
		{
			name: "abc < acba == true",
			args: args{
				x: "abc",
				y: "acba",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isLessString(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("isLessString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inventory_GetGroupsByPage(t *testing.T) {
	inv := NewInventory(make(map[string]Group))
	for i := 99; i >= 0; i-- {
		err := inv.AddHostToGroup(fmt.Sprintf("testGroup%d", i), ParseHost(fmt.Sprintf("192.168.%d.1", i)), ParseHost(fmt.Sprintf("10.10.%d.[1:254]", i)), ParseHost(fmt.Sprintf("172.%d.[1:100].[128:254]", i)))
		if err != nil {
			t.Fatal(err)
		}
	}
	type fields struct {
		sortedGroupNames []string
		groups           map[string]Group
		isTruncatedGroup map[string]bool
	}
	type args struct {
		recordsPerPage uint
		page           uint
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		wantErr            bool
		wantTotalGroupsNum uint
		wantTotalPagesNum  uint
		wantGroupsMap      map[string]Group
	}{
		{
			name: "normal test",
			fields: fields{
				sortedGroupNames: inv.(*inventory).sortedGroupNames,
				groups:           inv.(*inventory).groups,
				isTruncatedGroup: inv.(*inventory).isTruncatedGroup,
			},
			args: args{
				recordsPerPage: 5,
				page:           1,
			},
			wantTotalGroupsNum: 100,
			wantTotalPagesNum:  21,
			wantGroupsMap:      map[string]Group{"all": inv.GenerateGroupAll(), "testGroup0": inv.GetAllGroups()["testGroup0"], "testGroup1": inv.GetAllGroups()["testGroup1"], "testGroup2": inv.GetAllGroups()["testGroup2"], "testGroup3": inv.GetAllGroups()["testGroup3"]},
		},
		{
			name: "normal test page 2",
			fields: fields{
				sortedGroupNames: inv.(*inventory).sortedGroupNames,
				groups:           inv.(*inventory).groups,
				isTruncatedGroup: inv.(*inventory).isTruncatedGroup,
			},
			args: args{
				recordsPerPage: 5,
				page:           2,
			},
			wantTotalGroupsNum: 100,
			wantTotalPagesNum:  21,
			wantGroupsMap:      map[string]Group{"testGroup4": inv.GetAllGroups()["testGroup4"], "testGroup5": inv.GetAllGroups()["testGroup5"], "testGroup6": inv.GetAllGroups()["testGroup6"], "testGroup7": inv.GetAllGroups()["testGroup7"], "testGroup8": inv.GetAllGroups()["testGroup8"]},
		},
		{
			name: "normal test last page",
			fields: fields{
				sortedGroupNames: inv.(*inventory).sortedGroupNames,
				groups:           inv.(*inventory).groups,
				isTruncatedGroup: inv.(*inventory).isTruncatedGroup,
			},
			args: args{
				recordsPerPage: 5,
				page:           21,
			},
			wantTotalGroupsNum: 100,
			wantTotalPagesNum:  21,
			wantGroupsMap:      map[string]Group{"testGroup99": inv.GetAllGroups()["testGroup99"]},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := inventory{
				sortedGroupNames: tt.fields.sortedGroupNames,
				groups:           tt.fields.groups,
				isTruncatedGroup: tt.fields.isTruncatedGroup,
			}
			want := &Groups{
				TotalGroupsNum: &tt.wantTotalGroupsNum,
				TotalPagesNum:  &tt.wantTotalPagesNum,
				GroupsMap:      tt.wantGroupsMap,
			}
			got, err := i.GetGroupsByPage(tt.args.page, tt.args.recordsPerPage)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGroupsByPage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("GetGroupsByPage() got = %+v, want %+v", got, want)
			}
		})
	}
}
