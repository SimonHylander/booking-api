package transport

import (
	"github.com/simonhylander/booker"
)

// User model response
// swagger:response userResp
type swaggUserResponse struct {
	// in:body
	Body struct {
		*booker.User
	}
}

// Users model response
// swagger:response userListResp
type swaggUserListResponse struct {
	// in:body
	Body struct {
		Users []booker.User `json:"users"`
		Page  int          `json:"page"`
	}
}