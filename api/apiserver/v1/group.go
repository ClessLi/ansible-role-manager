package v1

type Group struct {
	GroupName string  `json:"group_name"`
	Hosts     []*Host `json:"hosts,omitempty"`
}

type Groups struct {
	TotalGroupsNum uint              `json:"totalGroupsNum"`
	TotalPagesNum  uint              `json:"totalPagesNum"`
	Items          map[string]*Group `json:"items"`
}
