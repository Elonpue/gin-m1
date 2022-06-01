package model

import "gorm.io/gorm"

type Menus struct {
	gorm.Model
	// Id       uint     `gorm:"primary_key" json:"id" `
	AuthName string   `gorm:"type:varchar(20)" json:"authName"`
	Path     string   `gorm:"type:varchar(20)" json:"path"`
	Sort     uint     `gorm:"tinyint(4)" json:"order" `
	ParentId *uint    `gorm:"default:0;comment:'父菜单编号(编号为0时表示根菜单)'" json:"parentId"`
	Children []*Menus `gorm:"-" json:"children"`
	Roles    []*Role  `gorm:"many2many:role_menus;" json:"roles"` // 角色菜单多对多关系
	// CreatedAt utils.Time `gorm:"type:timestamp" json:"created_at" `
	// UpdatedAt utils.Time `gorm:"type:timestamp" json:"updated_at" `
}
