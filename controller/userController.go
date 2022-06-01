package controller

import (
	"gin-m1/dto"
	"gin-m1/model"
	"gin-m1/response"
	"gin-m1/service"
	"gin-m1/utils"
	"gin-m1/vo"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
)

type IUserController interface {
	Login(c *gin.Context)
	GetUserInfo(c *gin.Context)
	GetUsers(c *gin.Context)
	ChangePwd(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type UserController struct {
	UserServe service.IUserService
}

func NewUserController() IUserController {
	userService := service.NewUserService()
	userController := UserController{UserServe: userService}
	return userController
}

func (uc UserController) Login(c *gin.Context) {
	var req vo.LoginReq

	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	u := &model.User{
		Username: req.Username,
		Password: req.Password,
	}

	user, err := uc.UserServe.Login(u)
	if err != nil {
		response.Fail(c, nil, "找不到用户")
		return
	}

	err = utils.ComparePasswd(user.Password, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"data": nil, "meta": gin.H{"status": 400, "msg": "密码错误"}})
		return
	}
	token, err := utils.ReleaseToken(user.ID)
	if err != nil {
		response.Fail(c, nil, "系统异常")
		log.Printf("token generate error: %v", err)
		return
	}

	// response.Success(c, gin.H{"token": token}, "登录成功")
	c.JSON(200, gin.H{"data": gin.H{"id": user.ID, "username": user.Username, "email": user.Email, "mobile": user.Mobile, "token": "Bearer " + token}, "meta": gin.H{"status": 200, "msg": "登录成功"}})

}

func (uc UserController) GetUserInfo(c *gin.Context) {
	user, err := uc.UserServe.GetsUserInfo(c)
	if err != nil {
		response.Fail(c, nil, "获取当前用户信息失败: "+err.Error())
		return
	}
	userInfoDto := dto.ToUserInfoDto(user)
	response.Success(c, gin.H{
		"userInfo": userInfoDto,
	}, "获取当前用户信息成功")
}

// 获取用户列表
func (uc UserController) GetUsers(c *gin.Context) {
	var req vo.UserListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// // 参数校验
	// if err := common.Validate.Struct(&req); err != nil {
	// 	errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
	// 	response.Fail(c, nil, errStr)
	// 	return
	// }

	// 获取
	users, total, pageNum, err := uc.UserServe.GetsUsers(&req)
	if err != nil {
		response.Fail(c, nil, "获取用户列表失败: "+err.Error())
		return
	}

	// response.Success(c, gin.H{"users": dto.ToUsersDto(users), "total": total}, "获取用户列表成功")
	c.JSON(200, gin.H{"data": gin.H{"users": dto.ToUsersDto(users), "totalpage": total, "pagenum": pageNum}, "meta": gin.H{"status": 200, "msg": "登录成功"}})
}

// 更新用户登录密码
func (uc UserController) ChangePwd(c *gin.Context) {
	var req vo.ChangePwdReq

	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	// 前端传来的密码是rsa加密的,先解密
	// 密码通过RSA解密
	req.OldPassword = utils.GenPasswd(req.OldPassword)

	req.NewPassword = utils.GenPasswd(req.NewPassword)

	// 获取当前用户
	user, err := uc.UserServe.GetsUserInfo(c)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 获取用户的真实正确密码
	correctPasswd := user.Password
	// 判断前端请求的密码是否等于真实密码
	err = utils.ComparePasswd(correctPasswd, req.OldPassword)
	if err != nil {
		response.Fail(c, nil, "原密码有误")
		return
	}
	// 更新密码
	err = uc.UserServe.ChangesPwd(user.Username, utils.GenPasswd(req.NewPassword))
	if err != nil {
		response.Fail(c, nil, "更新密码失败: "+err.Error())
		return
	}
	response.Success(c, nil, "更新密码成功")
}

// 创建用户
func (uc UserController) CreateUser(c *gin.Context) {
	var req vo.CreateUserReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	// 密码通过RSA解密
	// 密码不为空就解密
	if req.Password != "" {
		req.Password = utils.GenPasswd(req.Password)
		if len(req.Password) < 6 {
			response.Fail(c, nil, "密码长度至少为6位")
			return
		}
	}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	currentRoleSortMin, ctxUser, err := uc.UserServe.GetsUserMinRoleSort(c)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	// 获取前端传来的用户角色id
	reqRoleIds := []uint{2}
	// 根据角色id获取角色
	rr := service.NewRoleService()
	roles, err := rr.GetRolesByIds(reqRoleIds)
	if err != nil {
		response.Fail(c, nil, "根据角色ID获取角色信息失败: "+err.Error())
		return
	}
	if len(roles) == 0 {
		response.Fail(c, nil, "未获取到角色信息")
		return
	}
	var reqRoleSorts []int
	for _, role := range roles {
		reqRoleSorts = append(reqRoleSorts, int(role.Sort))
	}
	// 前端传来用户角色排序最小值（最高等级角色）
	reqRoleSortMin := uint(funk.MinInt(reqRoleSorts))

	// 当前用户的角色排序最小值 需要小于 前端传来的角色排序最小值（用户不能创建比自己等级高的或者相同等级的用户）
	if currentRoleSortMin >= reqRoleSortMin {
		response.Fail(c, nil, "用户不能创建比自己等级高的或者相同等级的用户")
		return
	}

	// 密码为空就默认123456
	if req.Password == "" {
		req.Password = "123456"
	}
	user := model.User{
		Username: req.Username,
		Password: utils.GenPasswd(req.Password),
		Mobile:   req.Mobile,
		Email:    req.Email,
		Status:   req.Status,
		Creator:  ctxUser.Username,
		Roles:    roles,
	}

	err = uc.UserServe.CreatesUser(&user)
	if err != nil {
		// response.Fail(c, nil, "创建用户失败: "+err.Error())
		c.JSON(200, gin.H{"meta": gin.H{"status": 201, "msg": "创建用户失败"}})
		return
	}
	// response.Success(c, nil, "创建用户成功")
	c.JSON(200, gin.H{"meta": gin.H{"status": 200, "msg": "创建用户成功"}})

}

// 更新用户
func (uc UserController) UpdateUser(c *gin.Context) {
	var req vo.CreateUserReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	//获取path中的userId
	userId, _ := strconv.Atoi(c.Param("userId"))
	if userId <= 0 {
		response.Fail(c, nil, "用户ID不正确")
		return
	}

	// 根据path中的userId获取用户信息
	oldUser, err := uc.UserServe.GetUserById(uint(userId))
	if err != nil {
		response.Fail(c, nil, "获取需要更新的用户信息失败: "+err.Error())
		return
	}

	// 获取当前用户
	ctxUser, err := uc.UserServe.GetsUserInfo(c)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 获取当前用户的所有角色
	currentRoles := ctxUser.Roles
	// 获取当前用户角色的排序，和前端传来的角色排序做比较
	var currentRoleSorts []int
	// 当前用户角色ID集合
	var currentRoleIds []uint
	for _, role := range currentRoles {
		currentRoleSorts = append(currentRoleSorts, int(role.Sort))
		currentRoleIds = append(currentRoleIds, role.ID)
	}
	// 当前用户角色排序最小值（最高等级角色）
	currentRoleSortMin := funk.MinInt(currentRoleSorts)

	// 获取前端传来的用户角色id
	reqRoleIds := req.RoleIds
	// 根据角色id获取角色
	rr := service.NewRoleService()
	roles, err := rr.GetRolesByIds(reqRoleIds)
	if err != nil {
		response.Fail(c, nil, "根据角色ID获取角色信息失败: "+err.Error())
		return
	}
	if len(roles) == 0 {
		response.Fail(c, nil, "未获取到角色信息")
		return
	}
	var reqRoleSorts []int
	for _, role := range roles {
		reqRoleSorts = append(reqRoleSorts, int(role.Sort))
	}
	// 前端传来用户角色排序最小值（最高等级角色）
	reqRoleSortMin := funk.MinInt(reqRoleSorts)

	user := model.User{
		Model:    oldUser.Model,
		Username: req.Username,
		Password: oldUser.Password,
		Mobile:   req.Mobile,
		Email:    req.Email,
		Status:   req.Status,
		Creator:  ctxUser.Username,
		Roles:    roles,
	}
	// 判断是更新自己还是更新别人
	if userId == int(ctxUser.ID) {
		// 如果是更新自己
		// 不能禁用自己
		if req.Status == 2 {
			response.Fail(c, nil, "不能禁用自己")
			return
		}
		// 不能更改自己的角色
		reqDiff, currentDiff := funk.Difference(req.RoleIds, currentRoleIds)
		if len(reqDiff.([]uint)) > 0 || len(currentDiff.([]uint)) > 0 {
			response.Fail(c, nil, "不能更改自己的角色")
			return
		}

		// 不能更新自己的密码，只能在个人中心更新
		if req.Password != "" {
			response.Fail(c, nil, "请到个人中心更新自身密码")
			return
		}

		// 密码赋值
		user.Password = ctxUser.Password

	} else {
		// 如果是更新别人
		// 用户不能更新比自己角色等级高的或者相同等级的用户
		// 根据path中的userIdID获取用户角色排序最小值
		minRoleSorts, err := uc.UserServe.GetsMinRoleSortsByIds([]uint{uint(userId)})
		if err != nil || len(minRoleSorts) == 0 {
			response.Fail(c, nil, "根据用户ID获取用户角色排序最小值失败")
			return
		}
		if currentRoleSortMin >= minRoleSorts[0] {
			response.Fail(c, nil, "用户不能更新比自己角色等级高的或者相同等级的用户")
			return
		}

		// 用户不能把别的用户角色等级更新得比自己高或相等
		if currentRoleSortMin >= reqRoleSortMin {
			response.Fail(c, nil, "用户不能把别的用户角色等级更新得比自己高或相等")
			return
		}

		// 密码赋值
		if req.Password != "" {

			user.Password = utils.GenPasswd(req.Password)
		}

	}

	// 更新用户
	err = uc.UserServe.UpdatesUser(&user)
	if err != nil {
		response.Fail(c, nil, "更新用户失败: "+err.Error())
		return
	}
	response.Success(c, nil, "更新用户成功")

}

// 批量删除用户
func (uc UserController) DeleteUser(c *gin.Context) {
	var req vo.DeleteUserReq
	// 参数绑定
	// if err := c.ShouldBind(&req); err != nil {
	// 	response.Fail(c, nil, err.Error())
	// 	return
	// }
	id := c.Param("id")
	u64, _ := strconv.ParseUint(id, 10, 32)
	id1 := uint(u64)
	req.UserIds = []uint{id1}
	// 前端传来的用户ID
	reqUserIds := req.UserIds
	// 根据用户ID获取用户角色排序最小值

	roleMinSortList, err := uc.UserServe.GetsMinRoleSortsByIds(reqUserIds)
	if err != nil || len(roleMinSortList) == 0 {
		response.Fail(c, nil, "根据用户ID获取用户角色排序最小值失败")
		return
	}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	minSort, ctxUser, err := uc.UserServe.GetsUserMinRoleSort(c)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	currentRoleSortMin := int(minSort)

	// 不能删除自己
	if funk.Contains(reqUserIds, ctxUser.ID) {
		response.Fail(c, nil, "用户不能删除自己")
		return
	}

	// 不能删除比自己角色排序低(等级高)的用户
	for _, sort := range roleMinSortList {
		if currentRoleSortMin >= sort {
			response.Fail(c, nil, "用户不能删除比自己角色等级高的用户")
			return
		}
	}

	err = uc.UserServe.DeleteUser(reqUserIds)
	if err != nil {
		// response.Fail(c, nil, "删除用户失败: "+err.Error())
		c.JSON(200, gin.H{"meta": gin.H{"status": 201, "msg": "删除用户失败"}})
		return
	}

	// response.Success(c, nil, "删除用户成功")
	c.JSON(200, gin.H{"meta": gin.H{"status": 200, "msg": "删除用户成功"}})

}
