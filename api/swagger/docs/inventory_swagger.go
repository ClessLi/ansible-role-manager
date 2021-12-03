package docs

import (
	v1 "github.com/ClessLi/ansible-role-manager/api/apiserver/v1"
	"github.com/ClessLi/ansible-role-manager/internal/pkg/core"
	metav1 "github.com/ClessLi/ansible-role-manager/internal/pkg/meta/v1"
)

// swagger:route POST /inventory/groups groups createGroupRequest
//
// Create groups.
//
// Create groups according to input parameters.
//
//     Security:
//       api_key:
//
//     Responses:
//       default: errResponse
//       200: createGroupResponse

// swagger:route DELETE /inventory/groups/{group} groups deleteGroupRequest
//
// Delete group.
//
// Delete group according to input parameters.
//
//     Responses:
//       default: errResponse
//       200: okResponse

// swagger:route PUT /inventory/groups/{group} groups updateGroupRequest
//
// Update group.
//
// Update group according to input parameters.
//
//     Security:
//       api_key:
//
//     Responses:
//       default: errResponse
//       200: updateGroupResponse

// swagger:route GET /inventory/groups/{group} groups getGroupRequest
//
// Get details for specified group.
//
// Get details for specified group according to input parameters.
//
//     Responses:
//       default: errResponse
//       200: getGroupResponse

// swagger:route GET /inventory/groups groups listGroupRequest
//
// List groups.
//
// List groups.
//
//     Responses:
//       default: errResponse
//       200: listGroupResponse

// List groups request.
// swagger:parameters listGroupRequest
type listGroupRequestParamsWrapper struct {
	// in:query
	metav1.ListOptions
}

// List groups response.
// swagger:response listGroupResponse
type listGroupResponseWrapper struct {
	// in:body
	Body v1.Groups
}

// Group response.
// swagger:response createGroupResponse
type createGroupResponseWrapper struct {
	// in:body
	Body v1.Group
}

// Group response.
// swagger:response updateGroupResponse
type updateGroupResponseWrapper struct {
	// in:body
	Body v1.Group
}

// Group response.
// swagger:response getGroupResponse
type getGroupResponseWrapper struct {
	// in:body
	Body v1.Group
}

// swagger:parameters createGroupRequest updateGroupRequest
type groupRequestParamsWrapper struct {
	// Group information.
	// in:body
	Body v1.Group
}

// swagger:parameters deleteGroupRequest getGroupRequest updateGroupRequest
type groupNameParamsWrapper struct {
	// Group name.
	// in:path
	Group string `json:"group"`
}

// ErrResponse defines the return message when an error occurred.
// swagger:response errResponse
type errResponseWrapper struct {
	// in:body
	Body core.ErrResponse
}

// Return nil json object.
// swagger:response okResponse
type okResponseWrapper struct{}
