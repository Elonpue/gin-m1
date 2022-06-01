package db

import (
	"fmt"
	"gin-m1/model"
	"gin-m1/utils"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func init() {
	utils.InitConfig()
	var err error
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Local",
		username,
		password,
		host,
		port,
		database)
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic(err.Error())
	}

	logger := utils.Log()
	logger.Info("mysql connect success")
	// err = Db.AutoMigrate(&model.User{})
	// if err != nil {
	// 	panic(err)
	// }
	// err1 := Db.AutoMigrate(&model.Role{})
	// if err1 != nil {
	// 	panic(err1)
	// }
	dbAutoMigrate()

}

// 自动迁移表结构
func dbAutoMigrate() {
	Db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Menus{},
	)
}

func GetDB() *gorm.DB {
	return Db
}
