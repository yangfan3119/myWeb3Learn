package controllers

import (
	"go01_4/config"
	"go01_4/databases"
	"go01_4/logger"
	"go01_4/models"
	"go01_4/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var mlog = logger.GetLogger()
var mdb *gorm.DB

func Init_db() {
	// 初始化数据库连接
	mdb = databases.GetDB()
	if mdb == nil {
		mlog.Error("数据库连接未初始化！")
	}
}

type UserCtr struct {
	mCfg config.BlogConfig
}

func NewUserCtr(cfg *config.BlogConfig) *UserCtr {
	return &UserCtr{mCfg: *cfg}
}

// 注册
func (c *UserCtr) Register(ctx *gin.Context) {
	// 从POST data Body中获取username、password、email
	var u models.User
	if err := ctx.ShouldBindBodyWithJSON(&u); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "用户信息获取异常.")
		return
	}
	// 验证用户名是否重复
	var users = models.User{}
	if err := mdb.Where("username=?", u.Username).First(&users).Error; err == nil {
		utils.Error(ctx, http.StatusBadRequest, "用户名或密码错误")
		return
	}
	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		mlog.Error("注册失败,密码加密失败。", err)
		utils.Error(ctx, http.StatusInternalServerError, "账号或密码错误")
		return
	}
	u.Password = string(hashedPassword)
	mlog.Infoln("注册用户:", u.Username, "密码：", u.Password, "邮箱：", u.Email)

	if err := mdb.Create(&u).Error; err != nil {
		utils.Error(ctx, http.StatusBadRequest, "账户注册失败")
		return
	}

	utils.Succ(ctx, "账户注册成功", nil)
}

// 登陆
func (c *UserCtr) Login(ctx *gin.Context) {
	// 使用用户名和密码登录，从POST data Body中获取username、password
	var getUser struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindBodyWithJSON(&getUser); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "用户信息获取异常.")
		return
	}
	// 验证用户名是否存在
	var dbUser = models.User{}
	if err := mdb.Where("username=?", getUser.Username).First(&dbUser).Error; err != nil {
		utils.Error(ctx, http.StatusBadRequest, "用户名或密码错误")
		return
	}
	// 验证密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(getUser.Password)); err != nil {
		utils.Error(ctx, http.StatusUnauthorized, "用户名或密码错误")
		return
	}
	// 生成JWT
	tokenString, err := utils.GenerateToken(getUser.Username, getUser.Password, dbUser.ID)
	if err != nil {
		utils.Error(ctx, http.StatusUnauthorized, "Token生成异常")
		return
	}
	utils.Succ(ctx, "登录成功", gin.H{
		"token":     tokenString,
		"user_id":   dbUser.ID,
		"user_name": dbUser.Username,
	})
}
