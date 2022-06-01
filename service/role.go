package service

import (
	"fmt"
	. "gin-m1/db"
	"gin-m1/model"
	"gin-m1/vo"
	"strings"
)

type IRoleService interface {
	GetRoles(req *vo.RoleListReq) ([]model.Role, int64, error) // 获取角色列表
	GetRolesByIds(roleIds []uint) ([]*model.Role, error)       // 根据角色ID获取角色
	CreateRole(role *model.Role) error                         // 创建角色
	UpdateRoleById(roleId uint, role *model.Role) error        // 更新角色
	BatchDeleteRoleByIds(roleIds []uint) error                 // 删除角色
}

type RoleService struct {
}

func NewRoleService() IRoleService {
	return RoleService{}
}

// 获取角色列表
func (r RoleService) GetRoles(req *vo.RoleListReq) ([]model.Role, int64, error) {
	var list []model.Role
	db := Db.Model(&model.Role{}).Order("created_at DESC")

	name := strings.TrimSpace(req.Name)
	if name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	keyword := strings.TrimSpace(req.Keyword)
	if keyword != "" {
		db = db.Where("keyword LIKE ?", fmt.Sprintf("%%%s%%", keyword))
	}
	status := req.Status
	if status != 0 {
		db = db.Where("status = ?", status)
	}
	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := int(req.PageNum)
	pageSize := int(req.PageSize)
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	return list, total, err
}

//根据角色ID获取角色
func (r RoleService) GetRolesByIds(roleIds []uint) ([]*model.Role, error) {
	var list []*model.Role
	err := Db.Where("id IN (?)", roleIds).Find(&list).Error
	return list, err
}

// 创建角色
func (r RoleService) CreateRole(role *model.Role) error {
	err := Db.Create(role).Error
	return err
}

// 更新角色
func (r RoleService) UpdateRoleById(roleId uint, role *model.Role) error {
	err := Db.Model(&model.Role{}).Where("id = ?", roleId).Updates(role).Error
	return err
}

// 删除角色
func (r RoleService) BatchDeleteRoleByIds(roleIds []uint) error {
	var roles []*model.Role
	err := Db.Where("id IN (?)", roleIds).Find(&roles).Error
	if err != nil {
		return err
	}
	err = Db.Select("Users", "Menus").Unscoped().Delete(&roles).Error

	return err
}
