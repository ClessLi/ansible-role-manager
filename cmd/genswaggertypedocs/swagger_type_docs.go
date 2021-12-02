package main

//go:generate swagger generate spec -o ../../api/swagger/swagger.yaml --scan-models

import (
	_ "github.com/ClessLi/ansible-role-manager/api/swagger/docs"
)

func main() {
}
