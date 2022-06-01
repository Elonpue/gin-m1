package controller

import (
	"gin-m1/model"
	"gin-m1/response"
	"gin-m1/service"
	"gin-m1/vo"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IRoleController interface {
	GetRoles(c *gin.Context)             // 获取角色列表
	CreateRole(c *gin.Context)           // 创建角色
	UpdateRoleById(c *gin.Context)       // 更新角色
	BatchDeleteRoleByIds(c *gin.Context) // 批量删除角色
}

type RoleController struct {
	RoleService service.IRoleService
}

func NewRoleController() IRoleController {
	roleService := service.NewRoleService()
	roleController := RoleController{RoleService: roleService}
	return roleController
}

// 获取角色列表
func (rc RoleController) GetRoles(c *gin.Context) {
	var req vo.RoleListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	// 获取角色列表
	roles, total, err := rc.RoleService.GetRoles(&req)
	if err != nil {
		response.Fail(c, nil, "获取角色列表失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"roles": roles, "total": total}, "获取角色列表成功")
}

// 创建角色
func (rc RoleController) CreateRole(c *gin.Context) {
	var req vo.CreateRoleReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 获取当前用户最高角色等级
	uc := service.NewUserService()
	sort, ctxUser, err := uc.GetsUserMinRoleSort(c)
	if err != nil {
		response.Fail(c, nil, "获取当前用户最高角色等级失败: "+err.Error())
		return
	}
	// 用户不能创建比自己等级高或相同等级的角色
	if sort >= req.Sort {
		response.Fail(c, nil, "不能创建比自己等级高或相同等级的角色")
		return
	}

	role := model.Role{
		Name:    req.Name,
		Keyword: req.Keyword,
		Desc:    &req.Desc,
		Status:  req.Status,
		Sort:    req.Sort,
		Creator: ctxUser.Username,
	}

	// 创建角色
	err = rc.RoleService.CreateRole(&role)
	if err != nil {
		response.Fail(c, nil, "创建角色失败: "+err.Error())
		return
	}
	response.Success(c, nil, "创建角色成功")

}

// 更新角色
func (rc RoleController) UpdateRoleById(c *gin.Context) {
	var req vo.CreateRoleReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	// 获取path中的roleId
	roleId, _ := strconv.Atoi(c.Param("roleId"))
	if roleId <= 0 {
		response.Fail(c, nil, "角色ID不正确")
		return
	}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	ur := service.NewUserService()
	minSort, ctxUser, err := ur.GetsUserMinRoleSort(c)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	// 不能更新比自己角色等级高或相等的角色
	// 根据path中的角色ID获取该角色信息
	roles, err := rc.RoleService.GetRolesByIds([]uint{uint(roleId)})
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	if len(roles) == 0 {
		response.Fail(c, nil, "未获取到角色信息")
		return
	}
	if minSort >= roles[0].Sort {
		response.Fail(c, nil, "不能更新比自己角色等级高或相等的角色")
		return
	}

	// 不能把角色等级更新得比当前用户的等级高
	if minSort >= req.Sort {
		response.Fail(c, nil, "不能把角色等级更新得比当前用户的等级高或相同")
		return
	}

	role := model.Role{
		Name:    req.Name,
		Keyword: req.Keyword,
		Desc:    &req.Desc,
		Status:  req.Status,
		Sort:    req.Sort,
		Creator: ctxUser.Username,
	}

	// 更新角色
	err = rc.RoleService.UpdateRoleById(uint(roleId), &role)
	if err != nil {
		response.Fail(c, nil, "更新角色失败: "+err.Error())
		return
	}

	// 如果更新成功，且更新了角色的keyword, 则更新casbin中policy

	// 更新角色成功处理用户信息缓存有两种做法:（这里使用第二种方法，因为一个角色下用户数量可能很多，第二种方法可以分散数据库压力）
	// 1.可以帮助用户更新拥有该角色的用户信息缓存,使用下面方法
	// err = ur.UpdateUserInfoCacheByRoleId(uint(roleId))
	// 2.直接清理缓存，让活跃的用户自己重新缓存最新用户信息
	ur.ClearUserInfoCache()

	response.Success(c, nil, "更新角色成功")
}

// 批量删除角色
func (rc RoleController) BatchDeleteRoleByIds(c *gin.Context) {
	var req vo.DeleteRoleReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	// 获取当前用户最高等级角色
	ur := service.NewUserService()
	minSort, _, err := ur.GetsUserMinRoleSort(c)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	// 前端传来需要删除的角色ID
	roleIds := req.RoleIds
	// 获取角色信息
	roles, err := rc.RoleService.GetRolesByIds(roleIds)
	if err != nil {
		response.Fail(c, nil, "获取角色信息失败: "+err.Error())
		return
	}
	if len(roles) == 0 {
		response.Fail(c, nil, "未获取到角色信息")
		return
	}

	// 不能删除比自己角色等级高或相等的角色
	for _, role := range roles {
		if minSort >= role.Sort {
			response.Fail(c, nil, "不能删除比自己角色等级高或相等的角色")
			return
		}
	}

	// 删除角色
	err = rc.RoleService.BatchDeleteRoleByIds(roleIds)
	if err != nil {
		response.Fail(c, nil, "删除角色失败")
		return
	}

	// 删除角色成功直接清理缓存，让活跃的用户自己重新缓存最新用户信息
	ur.ClearUserInfoCache()
	response.Success(c, nil, "删除角色成功")

}
