package handler

import (
	"fmt"
	"net/http"
	"sarkor/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"sarkor/internal/service"
)

type UserTokenData struct {
	UserID int
	Login  string
}

type UserHandler struct {
	svc service.UserSvc
}

func NewUserHandler(svc service.UserSvc) *UserHandler {
	return &UserHandler{svc: svc}
}

func (uh *UserHandler) AuthFreeRoutes(engine *gin.RouterGroup) {
	authFree := engine.Group("/user")
	{
		authFree.POST("/register", uh.Register)
		authFree.POST("/auth", uh.Auth)
	}
}

func (uh *UserHandler) AuthRoutes(engine *gin.RouterGroup) {
	auth := engine.Group("/user")
	{
		auth.GET("/:name", uh.GetUser)
		auth.POST("/phone", uh.AddPhone)
		auth.GET("/phone", uh.GetPhone)
		auth.PUT("/phone", uh.UpdatePhone)
		auth.DELETE("/phone/:phone", uh.DeletePhone)
	}
}

// Register
// @Param        login    		query   	string  true "login"
// @Param        password   	query   	string  true "password"
// @Param        username   	query   	string  true "username"
// @Param        age    		query   	int  true "age"
// @Success      200      		{object}	models.Res
// @Failure      400,500  		{object}	models.Res
// @Router       /user/register [post]
func (uh *UserHandler) Register(ctx *gin.Context) {
	var req models.UserRegReq
	var err error

	req.Login = ctx.PostForm("login")
	req.Password = ctx.PostForm("password")
	req.Name = ctx.PostForm("name")

	if req.Age, err = strconv.Atoi(ctx.PostForm("age")); err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, models.Res{Error: "Request form invalid"})
		return
	}

	if err = uh.svc.AddNewUser(ctx, req); err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, models.Res{Error: "Make sure using unique username and login"})
		return
	}

	ctx.JSON(http.StatusOK, models.Res{Success: true})
}

// Auth
// @Param        body  body   models.UserAuthReq true "body"
// @Success      200 {object} models.UserAuthRes
// @Failure 	 400 {object} models.UserAuthRes
// @Router       /user/auth [post]
func (uh *UserHandler) Auth(ctx *gin.Context) {
	var req models.UserAuthReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.UserAuthRes{Error: "Request body invalid"})
		return
	}

	token, err := uh.svc.Auth(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.UserAuthRes{Error: err.Error()})
		return
	}

	exp := time.Now().Add(365 * 24 * time.Hour)
	ctx.SetCookie("token", token, exp.Second(), "/", "", false, true)

	ctx.JSON(http.StatusOK, models.UserAuthRes{Token: token})
}

// GetUser
// @Param        name  		path   	 string true "name"
// @Success      200 		{object} models.UserRes
// @Failure 	 400,401 	{object} models.UserRes
// @Router       /user/{name} [get]
func (uh *UserHandler) GetUser(ctx *gin.Context) {
	userData, err := uh.svc.GetUserByName(ctx, ctx.Param("name"))
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, models.UserRes{Error: "No user found"})
		return
	}

	ctx.JSON(http.StatusOK, models.UserRes{Success: true, ID: userData.ID, Name: userData.Name, Age: userData.Age})
}

// AddPhone
// @Param        body  		body   	 models.UserPhoneReq true "body"
// @Success      200 		{object} models.Res
// @Failure 	 400,401 	{object} models.Res
// @Router       /user/phone [post]
func (uh *UserHandler) AddPhone(ctx *gin.Context) {
	tokenData := ctx.Request.Context().Value("user").(UserTokenData)
	if tokenData.UserID == 0 {
		ctx.JSON(http.StatusUnauthorized, models.Res{Error: "Authorization needed"})
		return
	}

	var req models.UserPhoneReq
	if err := ctx.BindJSON(&req); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, models.Res{Error: "Request body invalid"})
		return
	}

	if err := uh.svc.AddUserPhone(ctx, tokenData.UserID, req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Res{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, models.Res{Success: true})
}

// GetPhone
// @Param        q  		query    string true "query"
// @Success      200 		{object} models.PhonesRes
// @Failure 	 400,401 	{object} models.PhonesRes
// @Router       /user/phone [get]
func (uh *UserHandler) GetPhone(ctx *gin.Context) {
	phones, err := uh.svc.GetUserPhone(ctx, ctx.Query("q"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.PhonesRes{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, models.PhonesRes{Success: true, Data: phones})
}

// UpdatePhone
// @Param        body  		body     models.UserPhoneUpdateReq  true "body"
// @Success      200 		{object} models.Res
// @Failure 	 400,401 	{object} models.Res
// @Router       /user/phone [put]
func (uh *UserHandler) UpdatePhone(ctx *gin.Context) {
	tokenData := ctx.Request.Context().Value("user").(UserTokenData)
	if tokenData.UserID == 0 {
		ctx.JSON(http.StatusUnauthorized, models.Res{Error: "Authorization needed"})
		return
	}

	var req models.UserPhoneUpdateReq
	if err := ctx.BindJSON(&req); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, models.Res{Error: "Request body invalid"})
		return
	}

	if err := uh.svc.UpdateUserPhone(ctx, tokenData.UserID, req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Res{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, models.Res{Success: true})
}

// DeletePhone
// @Param        phone  	path     string  true "phone"
// @Success      200 		{object} models.Res
// @Failure 	 400,401 	{object} models.Res
// @Router       /user/phone/{phone} [delete]
func (uh *UserHandler) DeletePhone(ctx *gin.Context) {
	tokenData := ctx.Request.Context().Value("user").(UserTokenData)
	if tokenData.UserID == 0 {
		ctx.JSON(http.StatusUnauthorized, models.Res{Error: "Authorization needed"})
		return
	}

	fmt.Println(ctx.Param("phone"))
	if err := uh.svc.RemoveUserPhone(ctx, tokenData.UserID, ctx.Param("phone")); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Res{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, models.Res{Success: true})
}
