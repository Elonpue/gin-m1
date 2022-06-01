package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 返回前端
func Response(c *gin.Context, httpStatus int, code int, data gin.H, msg string) {
	c.JSON(httpStatus, gin.H{"code": code, "data": data, "msg": msg})
}

// 返回前端-成功
func Success(c *gin.Context, data gin.H, msg string) {
	Response(c, http.StatusOK, 200, data, msg)
}

// 返回前端-失败
func Fail(c *gin.Context, data gin.H, msg string) {
	Response(c, http.StatusBadRequest, 400, data, msg)
}