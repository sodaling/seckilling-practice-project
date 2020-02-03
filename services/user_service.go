package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"seckilling-practice-project/common"
	"seckilling-practice-project/models"
	"seckilling-practice-project/respsoiories"
)

type IUserService interface {
	IsPwdSuccess(userName string, pwd string) (user *models.User, isOk bool)
	AddUser(user *models.User) (userId int64, err error)
}

type UserService struct {
	UserRepository respsoiories.IUserRepository
}

func (u *UserService) IsPwdSuccess(userName string, pwd string) (user *models.User, isOk bool) {
	user, err := u.UserRepository.Select(userName)
	if err != nil {
		log.Panicln(err)
		return &models.User{}, false
	}
	isOk, _ = ValidatePassword(pwd, user.HashPassword)
	if !isOk {
		return &models.User{}, false
	}
	return user, true
}

func ValidatePassword(pwd string, hashed string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pwd)); err != nil {
		return false, errors.New("密码比对错误！")
	}
	return true, nil
}

func (u *UserService) AddUser(user *models.User) (userId int64, err error) {
	pwdByte, err := GeneratePassword(user.HashPassword)
	if err != nil {
		return 0, err
	}
	user.HashPassword = string(pwdByte)
	return u.UserRepository.Insert(user)
}

func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

func NewUserService(userRepository respsoiories.IUserRepository) IUserService {
	return &UserService{UserRepository: userRepository}
}

func DefaultUserSerivice() IUserService {
	mysqlCon, err := common.DefaultDb()
	if err != nil {
		panic(err)
	}
	userRepo := respsoiories.NewUserManagerRepository("user", mysqlCon)
	return NewUserService(userRepo)
}
