package service

import (
	"errors"
	"fmt"
	. "gin-m1/db"
	"gin-m1/model"
	"gin-m1/utils"
	"gin-m1/vo"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pmylund/go-cache"
	"github.com/thoas/go-funk"
)

type IUserService interface {
	Login(user *model.User) (*model.User, error) // 登录
	GetsUserInfo(c *gin.Context) (model.User, error)
	GetUserById(id uint) (model.User, error) // 获取单个用户
	GetsUsers(req *vo.UserListReq) ([]*model.User, int64, int, error)
	ChangesPwd(username string, newPasswd string) error
	CreatesUser(user *model.User) error
	UpdatesUser(user *model.User) error
	DeleteUser(ids []uint) error
	GetsUserMinRoleSort(c *gin.Context) (uint, model.User, error) // 获取当前用户角色排序最小值（最高等级角色）以及当前用户信息
	GetsMinRoleSortsByIds(ids []uint) ([]int, error)              // 根据用户ID获取用户角色排序最小值
	SetUserInfoCache(username string, user model.User)            // 设置用户信息缓存
	UpdateUserInfoCacheByRoleId(roleId uint) error                // 根据角色ID更新拥有该角色的用户信息缓存
	ClearUserInfoCache()                                          // 清理所有用户信息缓存
}

type UserService struct {
}

// 当前用户信息缓存，避免频繁获取数据库
var userInfoCache = cache.New(24*time.Hour, 48*time.Hour)

func NewUserService() IUserService {
	return UserService{}
}

// type User model.User

// //登陆
// func (u *User) Login() (b *User) {
// 	b = new(User)
// 	Db.Where("username = ? ", u.Username).First(b)
// 	return b
// }

// 登录
func (us UserService) Login(user *model.User) (*model.User, error) {
	// 根据用户名获取用户(正常状态:用户状态正常)
	var firstUser model.User
	err := Db.
		Where("username = ?", user.Username).
		Preload("Roles").
		First(&firstUser).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 判断用户的状态
	userStatus := firstUser.Status
	if userStatus != 1 {
		return nil, errors.New("用户被禁用")
	}

	// 判断用户拥有的所有角色的状态,全部角色都被禁用则不能登录
	roles := firstUser.Roles
	isValidate := false
	for _, role := range roles {
		// 有一个正常状态的角色就可以登录
		if role.Status == 1 {
			isValidate = true
			break
		}
	}

	if !isValidate {
		return nil, errors.New("用户角色被禁用")
	}

	// 校验密码
	err = utils.ComparePasswd(firstUser.Password, user.Password)
	if err != nil {
		return &firstUser, errors.New("密码错误")
	}
	return &firstUser, nil
}

// 获取当前登录用户信息
// 需要缓存，减少数据库访问
func (us UserService) GetsUserInfo(c *gin.Context) (model.User, error) {
	var NewUser model.User
	ctxUser, exist := c.Get("user")
	if !exist {
		return NewUser, errors.New("用户未登陆")
	}
	u, _ := ctxUser.(model.User)
	cacheUser, found := userInfoCache.Get(u.Username)
	var user model.User
	var err error
	if found {
		user = cacheUser.(model.User)
		err = nil
	} else {
		user, err = us.GetUserById(u.ID)
		if err != nil {
			userInfoCache.Delete(u.Username)
		} else {
			userInfoCache.Set(u.Username, user, cache.DefaultExpiration)
		}
	}
	return user, err
}

// 获取单个用户
func (us UserService) GetUserById(id uint) (model.User, error) {
	fmt.Println("GetUserById---")
	var user model.User
	err := Db.Where("id = ?", id).Preload("Roles").First(&user).Error
	return user, err
}

// 获取用户列表
func (us UserService) GetsUsers(req *vo.UserListReq) ([]*model.User, int64, int, error) {
	var list []*model.User
	db := Db.Model(&model.User{}).Order("created_at DESC")
	username := strings.TrimSpace(req.Username)
	if username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	mobile := strings.TrimSpace(req.Mobile)
	if mobile != "" {
		db = db.Where("mobile LIKE ?", fmt.Sprintf("%%%s%%", mobile))
	}
	email := strings.TrimSpace(req.Email)
	if email != "" {
		db = db.Where("email LIKE ?", fmt.Sprintf("%%%s%%", email))
	}

	status := req.Status
	if status != 0 {
		db = db.Where("status = ?", status)
	}
	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	var total int64
	var pageNum int
	err := db.Count(&total).Error
	if err != nil {
		return list, total, pageNum, err
	}
	pageNum = int(req.PageNum)
	pageSize := int(req.PageSize)
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Preload("Roles").Find(&list).Error
	} else {
		err = db.Preload("Roles").Find(&list).Error
	}
	return list, total, pageNum, err
}

// 获取当前用户角色排序最小值（最高等级角色）以及当前用户信息
func (us UserService) GetsUserMinRoleSort(c *gin.Context) (uint, model.User, error) {
	// 获取当前用户
	ctxUser, err := us.GetsUserInfo(c)
	if err != nil {
		return 999, ctxUser, err
	}
	// 获取当前用户的所有角色
	currentRoles := ctxUser.Roles
	// 获取当前用户角色的排序，和前端传来的角色排序做比较
	var currentRoleSorts []int
	for _, role := range currentRoles {
		currentRoleSorts = append(currentRoleSorts, int(role.Sort))
	}
	// 当前用户角色排序最小值（最高等级角色）
	currentRoleSortMin := uint(funk.MinInt(currentRoleSorts))

	return currentRoleSortMin, ctxUser, nil
}

func (us UserService) ChangesPwd(username string, hashNewPasswd string) error {

	err := Db.Model(&model.User{}).Where("username = ?", username).Update("password", hashNewPasswd).Error
	// 如果更新密码成功，则更新当前用户信息缓存
	// 先获取缓存
	cacheUser, found := userInfoCache.Get(username)
	if err == nil {
		if found {
			user := cacheUser.(model.User)
			user.Password = hashNewPasswd
			userInfoCache.Set(username, user, cache.DefaultExpiration)
		} else {
			// 没有缓存就获取用户信息缓存
			var user model.User
			Db.Where("username = ?", username).First(&user)
			userInfoCache.Set(username, user, cache.DefaultExpiration)
		}
	}

	return err
}

// 创建用户
func (us UserService) CreatesUser(user *model.User) error {
	err := Db.Create(user).Error
	return err
}

// 更新用户
func (us UserService) UpdatesUser(user *model.User) error {
	err := Db.Model(user).Updates(user).Error
	if err != nil {
		return err
	}
	err = Db.Model(user).Association("Roles").Replace(user.Roles)

	//err := common.DB.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&user).Error

	// 如果更新成功就更新用户信息缓存
	if err == nil {
		userInfoCache.Set(user.Username, *user, cache.DefaultExpiration)
	}
	return err
}

// 根据用户ID获取用户角色排序最小值
func (us UserService) GetsMinRoleSortsByIds(ids []uint) ([]int, error) {
	// 根据用户ID获取用户信息
	var userList []model.User
	err := Db.Where("id IN (?)", ids).Preload("Roles").Find(&userList).Error
	if err != nil {
		return []int{}, err
	}
	if len(userList) == 0 {
		return []int{}, errors.New("未获取到任何用户信息")
	}
	var roleMinSortList []int
	for _, user := range userList {
		roles := user.Roles
		var roleSortList []int
		for _, role := range roles {
			roleSortList = append(roleSortList, int(role.Sort))
		}
		roleMinSort := funk.MinInt(roleSortList)
		roleMinSortList = append(roleMinSortList, roleMinSort)
	}
	return roleMinSortList, nil
}

// 设置用户信息缓存
func (us UserService) SetUserInfoCache(username string, user model.User) {
	userInfoCache.Set(username, user, cache.DefaultExpiration)
}

// 根据角色ID更新拥有该角色的用户信息缓存
func (us UserService) UpdateUserInfoCacheByRoleId(roleId uint) error {

	var role model.Role
	err := Db.Where("id = ?", roleId).Preload("Users").First(&role).Error
	if err != nil {
		return errors.New("根据角色ID角色信息失败")
	}

	users := role.Users
	if len(users) == 0 {
		return errors.New("根据角色ID未获取到拥有该角色的用户")
	}

	// 更新用户信息缓存
	for _, user := range users {
		_, found := userInfoCache.Get(user.Username)
		if found {
			userInfoCache.Set(user.Username, *user, cache.DefaultExpiration)
		}
	}

	return err
}

// 批量删除
func (us UserService) DeleteUser(ids []uint) error {
	// 用户和角色存在多对多关联关系
	var users []model.User
	for _, id := range ids {
		// 根据ID获取用户
		user, err := us.GetUserById(id)
		if err != nil {
			return errors.New(fmt.Sprintf("未获取到ID为%d的用户", id))
		}
		users = append(users, user)
	}

	err := Db.Select("Roles").Unscoped().Delete(&users).Error
	// 删除用户成功，则删除用户信息缓存
	if err == nil {
		for _, user := range users {
			userInfoCache.Delete(user.Username)
		}
	}
	return err
}

//清理所有用户信息缓存
func (us UserService) ClearUserInfoCache() {
	userInfoCache.Flush()
}
