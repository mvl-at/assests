package assets

import "github.com/mvl-at/model"

type UserInfo struct {
	Member *model.Member `json:"members"`
	Roles  []model.Role  `json:"roles"`
}
