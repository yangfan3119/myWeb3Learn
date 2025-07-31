package controllers

import (
	"go01_4/middlewares"
	"go01_4/models"
	"go01_4/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommentCtr struct {
}

func NewCommentCtr() *CommentCtr {
	return &CommentCtr{}
}

func (c *CommentCtr) Create(ctx *gin.Context) {
	var newComm struct {
		Content string `json:"content"`
		PostID  string `json:"post_id"`
		UserID  string `json:"user_id"`
	}
	// 获取评论信息
	if err := ctx.ShouldBindBodyWithJSON(&newComm); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "评论数据异常，无法保存。")
		return
	}
	// 保存评论
	var comm = models.Comment{
		Content: newComm.Content,
		UserID:  utils.StrIdToUint(newComm.UserID),
		PostID:  utils.StrIdToUint(newComm.PostID),
	}
	if err := mdb.Create(&comm).Error; err != nil {
		utils.Error(ctx, http.StatusBadRequest, "评论数据异常，无法保存。")
		return
	}
	utils.Succ(ctx, "评论创建成功！", nil)
}

func (c *CommentCtr) Update(ctx *gin.Context) {
	// 评论权限确认，请求确权ID和评论ID是否一致。
	var newComm struct {
		CommentID string `json:"comment_id"`
		Content   string `json:"content"`
		PostID    string `json:"post_id"`
		UserID    string `json:"user_id"`
	}
	// 获取评论信息
	if err := ctx.ShouldBindBodyWithJSON(&newComm); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "评论数据异常，无法保存。")
		return
	}
	userId_str, _ := ctx.Get(middlewares.CtxUserIDKey)
	if userId_str != newComm.UserID || userId_str == "" {
		utils.Error(ctx, http.StatusBadRequest, "用户无权进行当前操作！")
		return
	}
	var comm models.Comment
	// 更新评论数据
	if err := mdb.
		Model(&comm).
		Where("comment_id=?", utils.StrIdToUint(newComm.CommentID)).
		Updates(map[string]interface{}{
			"Content": newComm.Content}).Error; err != nil {
		utils.Error(ctx, http.StatusBadRequest, "评论更新异常。")
		return
	}
	utils.Succ(ctx, "评论更新成功", nil)
}

func (c *CommentCtr) DelComm(ctx *gin.Context) {
	// 评论权限确认，请求确权ID和评论ID是否一致。
	var newComm struct {
		CommentID string `json:"comment_id"`
		Content   string `json:"content"`
		PostID    string `json:"post_id"`
		UserID    string `json:"user_id"`
	}
	// 获取评论信息
	if err := ctx.ShouldBindBodyWithJSON(&newComm); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "评论数据异常，无法删除。")
		return
	}
	userId_str, _ := ctx.Get(middlewares.CtxUserIDKey)
	if userId_str != newComm.UserID || userId_str == "" {
		utils.Error(ctx, http.StatusBadRequest, "用户无权进行当前操作！")
		return
	}
	var comm models.Comment
	// 更新评论数据
	if err := mdb.Where("comment_id=?", utils.StrIdToUint(newComm.CommentID)).Delete(&comm).Error; err != nil {
		utils.Error(ctx, http.StatusBadRequest, "评论删除异常。")
		return
	}
	utils.Succ(ctx, "评论删除成功", nil)
}

func (c *CommentCtr) GetComms(ctx *gin.Context) {
	postId := ctx.Param("id")
	if postId == "" {
		utils.Error(ctx, http.StatusBadRequest, "文章id获取异常")
		return
	}

	page := utils.StrIdToInt(ctx.DefaultQuery("page", "1"))
	limit := utils.StrIdToInt(ctx.DefaultQuery("limit", "30"))

	offIndex := (page - 1) * limit
	var comms []models.Comment
	if err := mdb.Where("post_id=?", utils.StrIdToUint(postId)).
		Order("create_at DESC").Offset(offIndex).Limit(limit).Find(&comms).Error; err != nil {
		utils.Error(ctx, http.StatusBadRequest, "评论获取异常")
		return
	}
	utils.Succ(ctx, "批次获取文章成功", comms)
}
