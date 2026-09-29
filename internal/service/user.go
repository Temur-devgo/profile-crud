package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"sarkor/internal/models"
	"sarkor/internal/repository"
)

type UserSvc struct {
	repo *repository.UsersRepo
}

func NewUserSvc(repo *repository.UsersRepo) *UserSvc {
	return &UserSvc{repo: repo}
}

func (us *UserSvc) AddNewUser(ctx *gin.Context, userData models.UserRegReq) error {
	password, err := bcrypt.GenerateFromPassword([]byte(userData.Password), 10)
	if err != nil {
		return err
	}

	if err = us.repo.AddNewUser(ctx, userData, password); err != nil {
		return err
	}

	return nil
}

type CustomClaims struct {
	Login  string `json:"login"`
	UserID int    `json:"user_id"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("my_secret_key")

func (us *UserSvc) Auth(ctx *gin.Context, req models.UserAuthReq) (string, error) {
	user, err := us.repo.GetUserByLogin(ctx, req.Login)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", fmt.Errorf("login %s not found", req.Login)
	}

	if ok := us.validatePassword(user.PassHash, req.Password); !ok {
		return "", fmt.Errorf("make sure login and password is correct")
	}

	cl := CustomClaims{Login: req.Login, UserID: user.ID}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cl)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (us *UserSvc) GetUserByName(ctx *gin.Context, userName string) (models.UserDb, error) {
	res, err := us.repo.GetUserByName(ctx, userName)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (us *UserSvc) AddUserPhone(ctx *gin.Context, userID int, phoneReq models.UserPhoneReq) error {
	if err := us.repo.AddUserPhone(ctx, userID, phoneReq); err != nil {
		return err
	}

	return nil
}

func (us *UserSvc) GetUserPhone(ctx *gin.Context, phone string) ([]models.UserPhoneDb, error) {
	res, err := us.repo.GetUserPhone(ctx, phone)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (us *UserSvc) UpdateUserPhone(ctx *gin.Context, userID int, req models.UserPhoneUpdateReq) error {
	if err := us.repo.UpdateUserPhone(ctx, userID, req); err != nil {
		return err
	}

	return nil
}

func (us *UserSvc) RemoveUserPhone(ctx *gin.Context, userID int, phone string) error {
	if err := us.repo.RemoveUserPhone(ctx, userID, phone); err != nil {
		return err
	}

	return nil
}

func (us *UserSvc) validatePassword(userPassHash, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(userPassHash), []byte(password)); err != nil {
		return false
	}

	return true
}
