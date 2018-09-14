package assets

import "github.com/mvl-at/model"

type UserInfo struct {
	Member *model.Member `json:"member"`
	Roles  []model.Role  `json:"roles"`
}
