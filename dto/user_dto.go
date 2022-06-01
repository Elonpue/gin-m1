package dto

import (
	"gin-m1/model"
	"time"
)

// 返回给前端的当前用户信息
type UserInfoDto struct {
	ID       uint          `json:"id"`
	Username string        `json:"username"`
	Mobile   string        `json:"mobile"`
	Email    string        `json:"email"`
	Roles    []*model.Role `json:"roles"`
}

func ToUserInfoDto(user model.User) UserInfoDto {
	return UserInfoDto{
		ID:       user.ID,
		Username: user.Username,
		Mobile:   user.Mobile,
		Email:    user.Email,
		Roles:    user.Roles,
	}
}

// 返回给前端的用户列表
type UsersDto struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Mobile    string    `json:"mobile"`
	Email     string    `json:"email"`
	Status    bool      `json:"mg_state"`
	Creator   string    `json:"creator"`
	CreatedAt time.Time `json:"create_time"`
	Role_name string    `json:"role_name"`
}

func ToUsersDto(userList []*model.User) []UsersDto {
	var users []UsersDto
	for _, user := range userList {
		userDto := UsersDto{
			ID:       user.ID,
			Username: user.Username,
			Mobile:   user.Mobile,
			Email:    user.Email,

			Creator:   user.Creator,
			CreatedAt: user.CreatedAt,
		}
		var staus bool
		if user.Status == 1 {
			staus = true
			userDto.Status = staus
		} else {
			staus = false
			userDto.Status = staus
		}

		var rolename string
		//  roles := make([]uint, 0)
		for _, role := range user.Roles {
			// roles = append(roles, role.ID)
			rolename = role.Name
		}
		userDto.Role_name = rolename
		users = append(users, userDto)
	}

	return users
}
