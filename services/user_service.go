package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"seckill_product/datamodels"
	"seckill_product/repositories"
)

// 接口
type IUserService interface {
	IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool)
	AddUser(user *datamodels.User) (userId int64, err error)
}

// 结构体
type UserService struct {
	UserRepository repositories.IUserRepository
}

// 构造函数
func NewService(repository repositories.IUserRepository) IUserService {
	return &UserService{UserRepository: repository}
}

// 接口方法

// 密码判断
func (u *UserService) IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool) {
	user, err := u.UserRepository.Select(userName)
	if err != nil {
		return
	}
	isOk, _ = ValidatePassWord(pwd, user.HashPassword)
	if !isOk {
		return &datamodels.User{}, false
	}
	return
}

// 添加用户
func (u *UserService) AddUser(user *datamodels.User) (userId int64, err error) {
	pwdByte, errPwd := GeneratePassWord(user.HashPassword)
	if errPwd != nil {
		return userId, errPwd
	}
	user.HashPassword = string(pwdByte)
	return u.UserRepository.Insert(user)
}

// 密码比对函数
func ValidatePassWord(userPassWord string, hashed string) (isOk bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassWord)); err != nil {
		return false, errors.New("密码比对错误！")
	}
	return true, nil
}

// 明文密码转密文密码
func GeneratePassWord(userPassWord string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassWord), bcrypt.DefaultCost)
}
