package controller

import (
	"github.com/gin-gonic/gin"

	"go-study/29-todo/middleware"
	"go-study/29-todo/model"
	"go-study/29-todo/utils"
)

type RegisterReq struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
}

type LoginReq struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.FailApp(c, utils.BadRequest("请提供 phone、password（至少 6 位）"))
		return
	}
	if req.Nickname == "" {
		req.Nickname = req.Phone
	}
	if _, err := model.GetUserByPhone(req.Phone); err == nil {
		utils.FailApp(c, utils.PhoneExists())
		return
	} else if !model.IsNotFound(err) {
		utils.FailApp(c, utils.WrapDB(err))
		return
	}
	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.FailApp(c, utils.Server("密码加密失败"))
		return
	}
	u := model.User{
		Name: req.Nickname, Username: req.Phone, Nickname: req.Nickname,
		Phone: req.Phone, Password: hashed, Role: "user",
	}
	if err := model.CreateUser(&u); err != nil {
		utils.FailApp(c, utils.WrapDB(err))
		return
	}
	utils.Created(c, "注册成功", gin.H{"user": gin.H{"id": u.ID, "phone": u.Phone, "nickname": u.Nickname}})
}

func Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.FailApp(c, utils.BadRequest("请提供 phone、password"))
		return
	}
	u, err := model.GetUserByPhone(req.Phone)
	if model.IsNotFound(err) {
		utils.FailApp(c, utils.UserNotFound())
		return
	}
	if err != nil {
		utils.FailApp(c, utils.WrapDB(err))
		return
	}
	if err := utils.CheckPassword(u.Password, req.Password); err != nil {
		utils.FailApp(c, utils.PasswordWrong())
		return
	}
	token, err := utils.GenerateToken(int64(u.ID), u.Role)
	if err != nil {
		utils.FailApp(c, utils.Server("生成 token 失败"))
		return
	}
	utils.OK(c, "登录成功", gin.H{
		"token": token, "token_type": "Bearer",
		"expires_in": int(utils.TokenTTL.Seconds()),
		"user":       gin.H{"id": u.ID, "phone": u.Phone, "nickname": u.Nickname},
	})
}

func Me(c *gin.Context) {
	u, err := model.GetUserByID(middleware.UID(c))
	if model.IsNotFound(err) {
		utils.FailApp(c, utils.UserNotFound())
		return
	}
	if err != nil {
		utils.FailApp(c, utils.WrapDB(err))
		return
	}
	utils.OK(c, "ok", gin.H{"user": gin.H{"id": u.ID, "phone": u.Phone, "nickname": u.Nickname}})
}
