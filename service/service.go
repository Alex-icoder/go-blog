package service

import "go-blog/dao"

// 登录服务接口
type Login interface {
	UserRegister(user dao.User)
	UserLogin(username string, password string)
}

// 文章管理服务
type Post interface {
	CreatePost(post dao.Post)
	FindPost(id uint) ([]dao.Post, error)
	UpdatePost(post dao.Post)
	deletePost(id uint)
}

// 评论管理服务
type Comment interface {
	CreateComment(comment dao.Comment)
	FindComment(id uint) ([]dao.Comment, error)
}
