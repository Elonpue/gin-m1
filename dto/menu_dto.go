package dto

import "gin-m1/model"

type MenuTreeDto struct {
	ID       uint          `json:"id"`
	AuthName string        `json:"authName"`
	Path     string        `json:"path"`
	Children []MenuTreeDto `json:"children"`
	Order    uint          `json:"order"`
}

func ToMenuTreeDto(menuList []*model.Menus) []MenuTreeDto {
	var menus []MenuTreeDto

	for _, m := range menuList {
		menuDto := MenuTreeDto{
			ID:       m.ID,
			AuthName: m.AuthName,
			Path:     m.Path,
			Order:    m.Sort,
		}
		children := make([]MenuTreeDto, 0)
		for _, c := range m.Children {

			m1 := MenuTreeDto{
				ID:       c.ID,
				AuthName: c.AuthName,
				Path:     c.Path,
				Order:    c.Sort,
			}
			children = append(children, m1)

		}
		menuDto.Children = children
		menus = append(menus, menuDto)
	}
	return menus
}
