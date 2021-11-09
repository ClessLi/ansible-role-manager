package ansible_inventory

import (
	"net"
	"reflect"
	"strings"
	"testing"
)

func Test_inventoryFileParser_Parse(t *testing.T) {
	invData := "[test-group]\n192.168.0.1\n10.1.0.1"
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Group
		wantErr bool
	}{
		{
			name: "normal inventory data",
			args: args{data: []byte(invData)},
			want: Group(&group{
				groupName: "test-group",
				hosts: []Host{
					Host(host{ipAddr: net.IPAddr{
						IP: net.IP{10, 1, 0, 1},
					}}),
					Host(host{ipAddr: net.IPAddr{
						IP: net.IP{192, 168, 0, 1},
					}}),
				},
			}),
		},
		{
			name:    "null data",
			args:    args{data: nil},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := parser{}
			got, err := i.Parse(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inventoryFileParser_Dump(t *testing.T) {
	testGroup := newGroup()
	err := testGroup.setName("test-group")
	if err != nil {
		t.Fatal(err)
	}
	_ = testGroup.AddHost(ParseHost("192.168.0.1"))
	_ = testGroup.AddHost(ParseHost("10.1.0.1"))
	data := []byte("[test-group]\n10.1.0.1\n192.168.0.1\n")
	type args struct {
		g Group
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "normal group input",
			args: args{g: testGroup},
			want: data,
		},
		{
			name:    "nil group",
			args:    args{g: nil},
			wantErr: true,
		},
		{
			name:    "null name group",
			args:    args{g: newGroup()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := parser{}
			got, err := i.Dump(tt.args.g)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dump() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.EqualFold(string(got), string(tt.want)) {
				t.Errorf("Dump() got = %s, want %s", got, tt.want)
			}
		})
	}
}
