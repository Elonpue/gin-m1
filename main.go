package main

import (
	"gin-m1/db"
	"gin-m1/router"
	"gin-m1/utils"

	"github.com/gin-gonic/gin"

	"github.com/spf13/viper"
)

func main() {
	//加载日志
	log := utils.Log()
	gin.SetMode(viper.GetString("server.model"))
	db.InitData()
	r := router.InitRouter()
	port := viper.GetString("server.port")
	if port != " " {
		log.Fatal(r.Run(":" + port))
	}

}
