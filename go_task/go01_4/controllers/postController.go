package controllers

import (
	"go01_4/middlewares"
	"go01_4/models"
	"go01_4/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostCtr struct {
}

func NewPostCtr() *PostCtr {
	return &PostCtr{}
}

// var mdb = databases.Mdb

func (p *PostCtr) Create(ctx *gin.Context) {
	var new_post models.Post
	if err := ctx.ShouldBindBodyWithJSON(&new_post); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "文章数据异常，无法保存。")
		return
	}
	if new_post.Title == "" || new_post.Content == "" {
		utils.Error(ctx, http.StatusBadRequest, "文章标题或内容不能为空")
		return
	}
	userId, exists := ctx.Get(middlewares.CtxUserIDKey)
	if !exists {
		utils.Error(ctx, http.StatusBadRequest, "userID获取异常")
		return
	}
	new_post.UserID = utils.StrIdToUint(userId.(string))
	if err := mdb.Create(&new_post).Error; err != nil {
		utils.Error(ctx, http.StatusBadRequest, "文章存储失败。")
		return
	}
	utils.Succ(ctx, "文章创建成功！", nil)
}

func (p *PostCtr) Update(ctx *gin.Context) {
	// 从地址中获取postId, 从上下文中获取UserId，进行所有权验证
	postId := ctx.Param("id")
	userId_str, _ := ctx.Get(middlewares.CtxUserIDKey)
	if postId == "" || userId_str == "" {
		utils.Error(ctx, http.StatusBadRequest, "文章id获取异常")
		return
	}
	var dbPost models.Post
	if err := mdb.Where("ID=?", postId).First(&dbPost).Error; err != nil {
		utils.Error(ctx, http.StatusBadRequest, "文章查询异常")
		return
	}
	userId := utils.StrIdToUint(userId_str.(string))
	if dbPost.UserID != userId {
		utils.Error(ctx, http.StatusBadRequest, "用户无权进行当前操作！")
		return
	}
	// 获取更新数据
	var newPost struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	if err := ctx.ShouldBindBodyWithJSON(&newPost); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "文章更新数据获取异常。")
		return
	}
	if err := mdb.Model(&dbPost).Updates(map[string]interface{}{
		"title":   newPost.Title,
		"content": newPost.Content}).Error; err != nil {
		utils.Error(ctx, http.StatusBadRequest, "文章更新异常。")
		return
	}
	utils.Succ(ctx, "文章更新成功", nil)
}

func (p *PostCtr) Delpost(ctx *gin.Context) {
	// 从地址中获取postId, 从上下文中获取UserId，进行所有权验证
	postId := ctx.Param("id")
	userId, _ := ctx.Get(middlewares.CtxUserIDKey)
	if postId == "" || userId == "" {
		utils.Error(ctx, http.StatusBadRequest, "文章id获取异常")
		return
	}
	var post models.Post
	if err := mdb.Where("ID=?", postId).First(&post).Error; err != nil {
		utils.Error(ctx, http.StatusBadRequest, "文章查询异常")
		return
	}
	if post.UserID != userId.(uint) {
		utils.Error(ctx, http.StatusBadRequest, "用户无权进行当前操作！")
		return
	}
	// 确权后，执行删除操作,此处需删除关联的comment 评论
	if err := mdb.Delete(&post).Error; err != nil {
		utils.Error(ctx, http.StatusBadRequest, "文章删除异常。")
		return
	}

	utils.Succ(ctx, "文章删除成功", nil)
}

func (p *PostCtr) GetPost(ctx *gin.Context) {
	// 从地址中获取postId
	postId := utils.StrIdToUint(ctx.Param("id"))
	// 从数据库中查询文章
	var dbpost models.Post
	if err := mdb.Where("id=?", postId).First(&dbpost).Error; err != nil {
		utils.Error(ctx, http.StatusBadRequest, "文章获取异常")
		return
	}
	utils.Succ(ctx, "评论获取成功", dbpost)
}

func (p *PostCtr) GetPostArray(ctx *gin.Context) {
	page := utils.StrIdToInt(ctx.DefaultQuery("page", "1"))
	limit := utils.StrIdToInt(ctx.DefaultQuery("limit", "30"))
	if page < 1 || limit < 1 {
		utils.Error(ctx, http.StatusBadRequest, "页码或每页数量参数错误")
		return
	}
	if limit > 100 {
		limit = 100 // 限制最大每页数量为100
	}

	offIndex := (page - 1) * limit
	var posts []models.Post
	if err := mdb.Order("created_at DESC").Offset(offIndex).Limit(limit).Find(&posts).Error; err != nil {
		utils.Error(ctx, http.StatusBadRequest, "文章获取异常")
		return
	}
	utils.Succ(ctx, "批次获取文章成功", posts)
}
