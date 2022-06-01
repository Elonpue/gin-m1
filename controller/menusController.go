package controller

import (
	"gin-m1/dto"
	"gin-m1/response"
	"gin-m1/service"

	"github.com/gin-gonic/gin"
)

func GetMenus(c *gin.Context) {
	var menu service.Menu

	result, err := menu.GetMenus()
	if err != nil {
		response.Fail(c, nil, "系统异常")
		return
	}

	res := dto.ToMenuTreeDto(result)

	// response.Success(c, gin.H{"data": result}, "获取成功")
	c.JSON(200, gin.H{"data": res, "meta": gin.H{"status": 200, "msg": "获取菜单列表成功"}})
}
