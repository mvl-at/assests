package assets

import "rest/model"

type UserInfo struct {
	Member *model.Member `json:"member"`
	Roles  []model.Role `json:"roles"`
}
