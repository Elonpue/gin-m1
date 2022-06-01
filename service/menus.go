package service

import (
	. "gin-m1/db"
	"gin-m1/model"
)

type Menu struct {
}

func (m *Menu) GetMenus() ([]*model.Menus, error) {
	var menus []*model.Menus
	err := Db.Order("`sort` asc").Find(&menus).Error
	return GenMenuTree(0, menus), err
}

func GenMenuTree(parentId uint, menus []*model.Menus) []*model.Menus {
	tree := make([]*model.Menus, 0)

	for _, m := range menus {
		if *m.ParentId == parentId {
			children := GenMenuTree(m.ID, menus)
			m.Children = children
			tree = append(tree, m)
		}
	}
	return tree
}
