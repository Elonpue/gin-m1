package db

import (
	"errors"
	"gin-m1/model"
	"gin-m1/utils"
	"log"

	"gorm.io/gorm"
)

func InitData() {
	newRoles := make([]*model.Role, 0)
	roles := []*model.Role{
		{
			Model:   gorm.Model{ID: 1},
			Name:    "管理员",
			Keyword: "admin",
			Desc:    new(string),
			Sort:    1,
			Status:  1,
			Creator: "系统",
		},
		{
			Model:   gorm.Model{ID: 2},
			Name:    "普通用户",
			Keyword: "user",
			Desc:    new(string),
			Sort:    3,
			Status:  1,
			Creator: "系统",
		},
		{
			Model:   gorm.Model{ID: 3},
			Name:    "访客",
			Keyword: "guest",
			Desc:    new(string),
			Sort:    5,
			Status:  1,
			Creator: "系统",
		},
	}
	for _, role := range roles {
		err := Db.First(&role, role.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newRoles = append(newRoles, role)
		}
	}
	if len(newRoles) > 0 {
		err := Db.Create(&newRoles).Error
		if err != nil {
			panic(err)
		}
	}
	// 写入菜单
	newMenus := make([]model.Menus, 0)
	var uint0 uint = 0
	var uint1 uint = 1
	var uint2 uint = 2
	var uint3 uint = 3
	var uint4 uint = 4
	var uint5 uint = 5
	menus := []model.Menus{
		{
			Model:    gorm.Model{ID: 1},
			AuthName: "用户管理",
			Path:     "users",
			Sort:     1,
			ParentId: &uint0,
			Roles:    roles[:1],
		},
		{
			Model:    gorm.Model{ID: 2},
			AuthName: "权限管理",
			Path:     "rights",
			Sort:     2,
			ParentId: &uint0,
			Roles:    roles[:1],
		},
		{
			Model:    gorm.Model{ID: 3},
			AuthName: "商品管理",
			Path:     "goods",
			Sort:     3,
			ParentId: &uint0,
			Roles:    roles[:1],
		},
		{
			Model:    gorm.Model{ID: 4},
			AuthName: "订单管理",
			Path:     "orders",
			Sort:     4,
			ParentId: &uint0,
			Roles:    roles[:1],
		},
		{
			Model:    gorm.Model{ID: 5},
			AuthName: "数据统计",
			Path:     "reports",
			Sort:     5,
			ParentId: &uint0,
			Roles:    roles[:1],
		},
		{
			Model:    gorm.Model{ID: 6},
			AuthName: "用户列表",
			Path:     "users",
			Sort:     6,
			ParentId: &uint1,
			Roles:    roles[:1],
		},
		{
			Model:    gorm.Model{ID: 7},
			AuthName: "角色列表",
			Path:     "roles",
			Sort:     7,
			ParentId: &uint2,
			Roles:    roles[:1],
		},
		{
			Model:    gorm.Model{ID: 8},
			AuthName: "权限列表",
			Path:     "rights",
			Sort:     8,
			ParentId: &uint2,
			Roles:    roles[:1],
		},
		{
			Model:    gorm.Model{ID: 9},
			AuthName: "商品列表",
			Path:     "goods",
			Sort:     9,
			ParentId: &uint3,
			Roles:    roles[:1],
		},
		{
			Model:    gorm.Model{ID: 10},
			AuthName: "分类参数",
			Path:     "params",
			Sort:     10,
			ParentId: &uint3,
			Roles:    roles[:1],
		},
		{
			Model:    gorm.Model{ID: 11},
			AuthName: "商品分类",
			Path:     "categories",
			Sort:     11,
			ParentId: &uint3,
			Roles:    roles[:1],
		},
		{
			Model:    gorm.Model{ID: 12},
			AuthName: "订单列表",
			Path:     "orders",
			Sort:     12,
			ParentId: &uint4,
			Roles:    roles[:1],
		},
		{
			Model:    gorm.Model{ID: 13},
			AuthName: "数据报表",
			Path:     "reports",
			Sort:     13,
			ParentId: &uint5,
			Roles:    roles[:1],
		},
	}
	for _, menu := range menus {
		err := Db.First(&menu, menu.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newMenus = append(newMenus, menu)
		}
	}
	if len(newMenus) > 0 {
		err := Db.Create(&newMenus).Error
		if err != nil {
			log.Fatal("写入系统菜单数据失败", err)
		}
	}
	// 写入用户
	newUsers := make([]model.User, 0)
	users := []model.User{
		{
			Model:    gorm.Model{ID: 1},
			Username: "admin",
			Password: utils.GenPasswd("123123"),
			Mobile:   "18888888888",
			Email:    "admin@qq.com",
			Status:   1,
			Creator:  "系统",
			Roles:    roles[:1],
		},
	}
	for _, user := range users {
		err := Db.First(&user, user.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newUsers = append(newUsers, user)
		}
	}

	if len(newUsers) > 0 {
		err := Db.Create(&newUsers).Error
		if err != nil {

			log.Fatal("写入用户数据失败", err)
		}
	}
}
