package vo

import "time"

// 用户登录结构体
type LoginReq struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// 获取用户列表结构体
type UserListReq struct {
	Username  string    `json:"username" form:"username" `
	Mobile    string    `json:"mobile" form:"mobile" `
	Email     string    `json:"email" form:"email" `
	Status    uint      `json:"status" form:"status" `
	CreatedAt time.Time `json:"created_at"`
	PageNum   uint      `json:"pageNum" form:"pageNum"`
	PageSize  uint      `json:"pageSize" form:"pageSize"`
}

// 更新密码结构体
type ChangePwdReq struct {
	OldPassword string `json:"oldPassword" form:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" form:"newPassword" validate:"required"`
}

// 创建用户结构体
type CreateUserReq struct {
	Username string `form:"username" json:"username" validate:"required,min=2,max=20"`
	Password string `form:"password" json:"password"`
	Mobile   string `form:"mobile" json:"mobile" validate:"required,checkMobile"`
	Email    string `form:"email" json:"email"`
	Status   uint   `form:"status" json:"status" validate:"oneof=1 2"`
	RoleIds  []uint `form:"roleIds" json:"roleIds" validate:"required"`
}

// 批量删除用户结构体
type DeleteUserReq struct {
	UserIds []uint `json:"userIds" form:"userIds"`
}
