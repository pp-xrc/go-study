package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"go-study/28-modular/middleware"
	"go-study/28-modular/model"
	"go-study/28-modular/utils"
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

type UserVO struct {
	ID       uint   `json:"id"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
}

func toVO(u model.User) UserVO {
	return UserVO{ID: u.ID, Phone: u.Phone, Nickname: u.Nickname, Role: u.Role}
}

// Register 校验参数 → 查重 → 加密 → 调 model.CreateUser
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
		Name:     req.Nickname,
		Username: req.Phone,
		Nickname: req.Nickname,
		Phone:    req.Phone,
		Password: hashed,
		Role:     "user",
	}
	if err := model.CreateUser(&u); err != nil {
		utils.FailApp(c, utils.WrapDB(err))
		return
	}
	utils.Created(c, "注册成功", gin.H{"user": toVO(u)})
}

// Login 校验 → GetUserByPhone → 比对密码 → 签发 JWT
func Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.FailApp(c, utils.BadRequest("请提供 phone、password"))
		return
	}

	u, err := model.GetUserByPhone(req.Phone)
	if model.IsNotFound(err) {
		utils.FailApp(c, utils.UserNotFound()) // 10001
		return
	}
	if err != nil {
		utils.FailApp(c, utils.WrapDB(err))
		return
	}

	if err := utils.CheckPassword(u.Password, req.Password); err != nil {
		utils.FailApp(c, utils.PasswordWrong()) // 10002
		return
	}

	token, err := utils.GenerateToken(int64(u.ID), u.Role)
	if err != nil {
		utils.FailApp(c, utils.Server("生成 token 失败"))
		return
	}
	utils.OK(c, "登录成功", gin.H{
		"token":      token,
		"token_type": "Bearer",
		"expires_in": int(utils.TokenTTL.Seconds()),
		"user":       toVO(u),
	})
}

func Me(c *gin.Context) {
	id := uint(c.GetInt64(middleware.CtxUserID))
	u, err := model.GetUserByID(id)
	if model.IsNotFound(err) {
		utils.FailApp(c, utils.UserNotFound())
		return
	}
	if err != nil {
		utils.FailApp(c, utils.WrapDB(err))
		return
	}
	utils.OK(c, "获取当前用户成功", gin.H{"user": toVO(u)})
}

func List(c *gin.Context) {
	list, err := model.GetUserList()
	if err != nil {
		utils.FailApp(c, utils.WrapDB(err))
		return
	}
	vos := make([]UserVO, 0, len(list))
	for _, u := range list {
		vos = append(vos, toVO(u))
	}
	utils.OK(c, "获取列表成功", gin.H{"count": len(vos), "list": vos})
}

func Detail(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id64 == 0 {
		utils.FailApp(c, utils.BadRequest("id 必须是正整数"))
		return
	}
	current := uint(c.GetInt64(middleware.CtxUserID))
	role := c.GetString(middleware.CtxRole)
	if role != "admin" && current != uint(id64) {
		utils.FailApp(c, utils.Forbidden("只能查看自己的信息"))
		return
	}

	u, err := model.GetUserByID(uint(id64))
	if model.IsNotFound(err) {
		utils.FailApp(c, utils.UserNotFound())
		return
	}
	if err != nil {
		utils.FailApp(c, utils.WrapDB(err))
		return
	}
	utils.OK(c, "获取详情成功", gin.H{"user": toVO(u)})
}
