package codes

import (
	"fmt"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/*
进阶gorm
题目1：模型定义
假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
要求 ：
使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章），
 Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
编写Go代码，使用Gorm创建这些模型对应的数据库表。
题目2：关联查询
基于上述博客系统的模型定义。
要求 ：
编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
编写Go代码，使用Gorm查询评论数量最多的文章信息。
题目3：钩子函数
继续使用博客系统的模型。
要求 ：
为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，
则更新文章的评论状态为 "无评论"。
*/
// 题目1,2,3
type User struct { // 用户
	gorm.Model
	UserCode  string `gorm:"size:16;not null;uniqueIndex;"`
	Name      string `gorm:"size:30;uniqueIndex;not null"`
	Age       uint8
	PostCount int       `gorm:"default:0"`
	Posts     []Post    `gorm:"foreignKey:UserCode;references:UserCode"`
	Comments  []Comment `gorm:"foreignKey:UserCode;references:UserCode"`
}

func (u *User) getPrintf() string {
	var ps, cms string
	for index, p := range u.Posts {
		ps = ps + fmt.Sprintf("posts %d Title: %s, Content: %s ", index, p.Title, p.Content)
	}
	for index, c := range u.Comments {
		cms = cms + fmt.Sprintf("comment %d, %s", index, c.Info)
	}
	return fmt.Sprintf("User name:%s, age:%d, posts:%s, comments:%s", u.Name, u.Age, ps, cms)
}

type Post struct { // 文章
	gorm.Model
	PostCode      string `gorm:"size:16;not null;uniqueIndex;"`
	Title         string
	Content       string
	CommentStatus string    `gorm:"default:'无评论'"`
	UserCode      string    `gorm:"not null"`
	Comments      []Comment `gorm:"foreignKey:PostCode;references:PostCode"`
}

type Comment struct { // 评论
	gorm.Model
	Info     string `gorm:"column:info"`
	UserCode string `gorm:"not null"`
	PostCode string `gorm:"not null"`
}

func (p *Post) AfterCreate(tx *gorm.DB) error {
	return tx.Transaction(func(tx *gorm.DB) error {
		var ux = &User{}
		if err := tx.Where("user_code", p.UserCode).First(&ux).Error; err != nil {
			return err
		}
		fmt.Println("Post 创建前 用户", ux.Name, "的文章数为：", ux.PostCount)

		if err := tx.Model(&User{}).Where("user_code = ?", p.UserCode).Update("post_count", ux.PostCount+1).Error; err != nil {
			fmt.Println("PostCount 更新出错，err:", err)
			return err
		}

		tx.Where("user_code", p.UserCode).First(&ux)
		fmt.Println("Post 创建后 用户", ux.Name, "的文章数为：", ux.PostCount)

		return nil
	})
}

func (c *Comment) AfterCreate(tx *gorm.DB) error {
	var count int64
	if err := tx.Model(&Comment{}).Where("post_code", c.PostCode).Count(&count).Error; err != nil {
		return err
	}
	cStatus := "无评论"
	if count > 0 {
		cStatus = strconv.FormatInt(count, 10)
	}
	if err := tx.Model(&Post{}).Where("post_code=?", c.PostCode).Update("comment_status", cStatus).Error; err != nil {
		return err
	}

	var px = &Post{}
	tx.Where("post_code=?", c.PostCode).First(&px)
	fmt.Println("Comment AfterCreate 《", px.Title, "》 comment_status:", px.CommentStatus)
	return nil
}

func (c *Comment) AfterDelete(tx *gorm.DB) error {
	var count int64
	var pCode = c.PostCode
	if err := tx.Model(&Comment{}).Where("post_code", c.PostCode).Count(&count).Error; err != nil {
		return err
	}
	cStatus := "无评论"
	if count > 0 {
		cStatus = strconv.FormatInt(count, 10)
	}
	if err := tx.Model(&Post{}).Where("post_code=?", c.PostCode).Update("comment_status", cStatus).Error; err != nil {
		return err
	}
	fmt.Println("Comment AfterDelete post_code:", pCode)
	var px = &Post{}
	tx.Where("post_code=?", c.PostCode).First(&px)
	fmt.Println("Comment AfterDelete 《", px.Title, "》 comment_status:", px.CommentStatus)

	return nil
}

func blog_create(db *gorm.DB) {
	var users = []User{
		{Name: "张三", Age: 20, UserCode: "u1"},
		{Name: "李四", Age: 19, UserCode: "u2"},
		{Name: "王五", Age: 14, UserCode: "u3"},
		{Name: "周大", Age: 18, UserCode: "u4"},
		{Name: "小六", Age: 13, UserCode: "u5"},
	}
	var posts = []Post{
		{Title: "Web3基础", Content: "WWWWWW", PostCode: "p1", UserCode: "u3"},
		{Title: "Web3进阶", Content: "EBEBEBEB", PostCode: "p2", UserCode: "u5"},
		{Title: "测试文章", Content: "EBEBEBEB", PostCode: "p4", UserCode: "u2"},
	}
	var comments = []Comment{
		{Info: "测试评论", UserCode: "u1", PostCode: "p4"},
		{Info: "文章太赞了", UserCode: "u2", PostCode: "p1"},
		{Info: "非常棒！棒棒的", UserCode: "u1", PostCode: "p1"},
		{Info: "文字太好了", UserCode: "u3", PostCode: "p1"},
		{Info: "不感兴趣。", UserCode: "u4", PostCode: "p1"},
		{Info: "牛文！！！", UserCode: "u1", PostCode: "p2"},
		{Info: "一般一般吧", UserCode: "u3", PostCode: "p2"},
		{Info: "文章太赞了", UserCode: "u2", PostCode: "p2"},
	}

	db.Migrator().DropTable(&User{}, &Post{}, &Comment{})
	db.AutoMigrate(&User{}, &Post{}, &Comment{})
	db.Create(users)
	db.Create(posts)
	db.Create(comments)
}

func Blog_Run() {
	var dsn string = "host=localhost user=postgres password=lyan123 dbname=gorm_test port=5432 sslmode=disable TimeZone=Asia/Shanghai "
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败: " + err.Error())
	} else {
		fmt.Println("数据库连接成功！")
	}
	db.Migrator().DropTable(&User{}, &Post{}, &Comment{})
	db.AutoMigrate(&User{}, &Post{}, &Comment{})
	blog_create(db)

	var user User
	if err := db.Preload("Posts").Preload("Comments").Where("name=?", "张三").Find(&user).Error; err != nil {
		fmt.Println("博客查询，张三信息查询失败")
		panic(err)
	}
	fmt.Println("博客查询，张三信息查询: ", user.getPrintf())

	type PostCommentCount struct {
		PostID       uint
		PostCode     string
		Title        string
		CommentCount int64
	}

	var commentCounts []PostCommentCount
	db.Model(&Post{}).
		Select("posts.id as post_id, posts.post_code, posts.title, COUNT(comments.id) as comment_count").
		Joins("LEFT JOIN comments ON comments.post_code = posts.post_code").
		Group("posts.id, posts.post_code, posts.title").
		Order("comment_count DESC").
		Scan(&commentCounts)

	for _, a := range commentCounts {
		fmt.Println("获取当前评论数情况：", a)
	}

	// 2. 找出最大评论数量
	var maxCount int64 = 0
	for _, count := range commentCounts {
		if count.CommentCount > maxCount {
			maxCount = count.CommentCount
		}
	}

	// 3. 查询评论数量等于最大值的文章
	var mostCommentedPosts []Post
	db.Model(&Post{}).
		Select("posts.*, COUNT(comments.id) as comment_count").
		Joins("LEFT JOIN comments ON comments.post_code = posts.post_code").
		Group("posts.id").
		Having("COUNT(comments.id) = ?", maxCount).
		Find(&mostCommentedPosts)

	// 打印结果
	fmt.Printf("\n评论数量最多的文章（评论数: %d）:\n", maxCount)
	for _, post := range mostCommentedPosts {
		fmt.Printf("文章ID: %d, 编码: %s, 标题: %s\n", post.ID, post.PostCode, post.Title)

		// 预加载该文章的评论
		var comments []Comment
		db.Where("post_code = ?", post.PostCode).Find(&comments)

		fmt.Println("评论列表:")
		for _, comment := range comments {
			fmt.Printf("  - %s (用户: %s)\n", comment.Info, comment.UserCode)
		}
		fmt.Println()
	}

	//为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
	var newP1 = &Post{Title: "Post基础", Content: "POSTPOSTPOST", PostCode: "p3", UserCode: "u1"}
	if err := db.Create(&newP1).Error; err != nil {
		panic(err)
	}

	// 为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
	// 测试文章为p4, 评论仅1条，comment id为1
	var p4 = &Post{}
	if err := db.Where("post_code=?", "p4").First(&p4).Error; err != nil {
		panic(err)
	}
	fmt.Println("删除测试，删除前文章的评论状态为：", p4.CommentStatus)
	var c1 = &Comment{}
	db.Where("id=?", 1).First(&c1)
	if err := db.Delete(&c1).Error; err != nil {
		panic(err)
	}
}
