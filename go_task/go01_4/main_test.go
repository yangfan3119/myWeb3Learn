package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go01_4/config"
	"go01_4/databases"
	"go01_4/logger"
	"go01_4/middlewares"
	"go01_4/models"
	"go01_4/router"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func setupTestApp() *gin.Engine {
	config.Cfg.SqliteName = ":memory:"
	config.Cfg.ServerPort = "8081"
	config.Cfg.JwtSecret = "myBlogSecret"

	logger.Init()

	databases.InitDB(config.Cfg.SqliteName)
	models.MigrateDB(databases.GetDB())

	db := databases.GetDB()
	// 添加3组用户、文章、评论测试数据
	users := []models.User{
		{Username: "user1", Password: "$2a$10$SJNHwr5mrmcD5S5g887Z3.OX3fhw4dnd.XhhUEaOR8wW3V1ot8aN.", Email: "user1@test.com"},
		{Username: "user2", Password: "$2a$10$/J4gz.aheRZc6vLI2HrgIOPyG7UqmN/ndWUBOJ4HnILOdjWSNuNNC", Email: "user2@test.com"},
		{Username: "user3", Password: "$2a$10$39OzxaUhP5Sw4GBligDpc.ebAhQc77shI4Cw5uIGI9tmIzDVxw.2y", Email: "user3@test.com"},
	}
	for i := range users {
		db.Create(&users[i])
	}
	posts := []models.Post{
		{Title: "Post1", Content: "Content1", UserID: 1},
		{Title: "Post2", Content: "Content2", UserID: 2},
		{Title: "Post3", Content: "Content3", UserID: 3},
	}
	for i := range posts {
		db.Create(&posts[i])
	}
	comments := []models.Comment{
		{Content: "Comment1", UserID: 1, PostID: 1},
		{Content: "Comment2", UserID: 2, PostID: 2},
		{Content: "Comment3", UserID: 3, PostID: 3},
	}
	for i := range comments {
		db.Create(&comments[i])
	}

	r := gin.New()
	r.Use(middlewares.RecoveryMiddleware())
	router.SetRouter(r)
	return r
}

// 注册
func TestRegister(t *testing.T) {
	app := setupTestApp()
	ts := httptest.NewServer(app)
	defer ts.Close()

	type test struct {
		input  string
		output int
	}
	tests := map[string]test{
		"No1": {input: `{"username":"user4","password":"pass4","email":"newuser@test.com"}`, output: http.StatusOK},
		"No2": {input: `{"username":"user1","password":"pass2","email":"newuser@test.com"}`, output: http.StatusBadRequest},
		"No3": {input: `{}`, output: http.StatusBadRequest},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := http.Post(ts.URL+"/register", "application/json", bytes.NewBufferString(tc.input))
			if err != nil {
				t.Fatalf("Register failed: %v", err)
			}
			defer resp.Body.Close()
			assert.Equal(t, tc.output, resp.StatusCode)
		})
	}
}

// 登录
func TestLogin(t *testing.T) {
	app := setupTestApp()
	ts := httptest.NewServer(app)
	defer ts.Close()

	type test struct {
		input  string
		output int
	}
	tests := map[string]test{
		"No1": {input: `{"username":"user1","password":"pass1"}`, output: http.StatusOK},
		"No2": {input: `{"username":"user1","password":"wrongpass"}`, output: http.StatusUnauthorized},
		"No3": {input: `{"username":"user1"}`, output: http.StatusUnauthorized},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := http.Post(ts.URL+"/login", "application/json", bytes.NewBufferString(tc.input))
			if err != nil {
				t.Fatalf("Register failed: %v", err)
			}
			defer resp.Body.Close()
			assert.Equal(t, tc.output, resp.StatusCode)
		})
	}
}

// 获取文章列表
func TestGetPostArray(t *testing.T) {
	app := setupTestApp()
	ts := httptest.NewServer(app)
	defer ts.Close()

	type test struct {
		input  string
		output int
	}
	tests := map[string]test{
		"No1": {input: "?page=1&limit=10", output: http.StatusOK},
		"No2": {input: "?page=0&limit=10", output: http.StatusBadRequest},
		"No3": {input: "?page=1&limit=0", output: http.StatusBadRequest},
		"No4": {input: "?page=1&limit=101", output: http.StatusOK}, // 限制最大每页数量为100
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := http.Get(ts.URL + "/post/" + tc.input)
			if err != nil {
				t.Fatalf("GetPostArray failed: %v", err)
			}
			defer resp.Body.Close()
			assert.Equal(t, tc.output, resp.StatusCode)
		})
	}
}

// 获取单篇文章
func TestGetPost(t *testing.T) {
	app := setupTestApp()
	ts := httptest.NewServer(app)
	defer ts.Close()

	type test struct {
		input  string
		output int
	}
	tests := map[string]test{
		"No1": {input: "1", output: http.StatusOK},
		"No2": {input: "999", output: http.StatusBadRequest}, // 假设不存在的文章ID
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := http.Get(ts.URL + "/post/" + tc.input)
			if err != nil {
				t.Fatalf("GetPost failed: %v", err)
			}
			defer resp.Body.Close()
			assert.Equal(t, tc.output, resp.StatusCode)
		})
	}
}

// 获取文章评论
func TestGetPostComments(t *testing.T) {
	app := setupTestApp()
	ts := httptest.NewServer(app)
	defer ts.Close()

	type test struct {
		input  string
		output int
	}
	tests := map[string]test{
		"No1": {input: "1/comment?page=1&limit=10", output: http.StatusOK},
		"No2": {input: "1/comment?page=0&limit=10", output: http.StatusBadRequest},
		"No3": {input: "1/comment?page=1&limit=0", output: http.StatusBadRequest},
		"No4": {input: "1/comment?page=1&limit=101", output: http.StatusOK},  // 限制最大每页数量为100
		"No5": {input: "999/comment?page=1&limit=10", output: http.StatusOK}, // 假设不存在的文章ID
		"No6": {input: "1/comment", output: http.StatusOK},                   // 缺少分页参数
		"No7": {input: "1/comment?page=1", output: http.StatusOK},            // 缺少limit参数
		"No8": {input: "1/comment?limit=10", output: http.StatusOK},          // 缺少page参数
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := http.Get(ts.URL + "/post/" + tc.input)
			if err != nil {
				t.Fatalf("GetPostComments failed: %v", err)
			}
			defer resp.Body.Close()
			assert.Equal(t, tc.output, resp.StatusCode)
		})
	}
}

// 需认证的接口测试（以创建文章为例）
func TestCreatePostAuth(t *testing.T) {
	app := setupTestApp()
	ts := httptest.NewServer(app)
	defer ts.Close()

	// 先登录获取token
	loginBody := `{"username":"user1","password":"pass1"}`
	resp, err := http.Post(ts.URL+"/Login", "application/json", bytes.NewBufferString(loginBody))
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	defer resp.Body.Close()
	var loginResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&loginResp)
	token := ""
	if data, ok := loginResp["data"].(map[string]interface{}); ok {
		token, _ = data["token"].(string)
	}
	if token == "" {
		t.Skip("No token, skip auth test")
	}

	// 创建文章
	postBody := `{"title":"new post","content":"new content","user_id":1}`
	req, _ := http.NewRequest("POST", ts.URL+"/post/create", bytes.NewBufferString(postBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("CreatePost failed: %v", err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK && resp2.StatusCode != http.StatusBadRequest {
		t.Errorf("CreatePost: expected 200 or 400, got %d", resp2.StatusCode)
	}
}
