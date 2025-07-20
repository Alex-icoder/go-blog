package service

import (
	"errors"
	"go-blog/config"
	"go-blog/dao"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
)

func UserRegister(user dao.User) error {
	var count int64
	config.DB.Model(&dao.User{}).Where("username = ? or email = ?", user.Username, user.Email).Count(&count)
	if count > 0 {
		return errors.New("当前用户名或邮箱已存在")
	}
	if err := config.DB.Create(&user).Error; err != nil {
		log.Println(err)
		return errors.New("注册用户失败")
	}
	return nil
}

func UserLogin(username string, password string) (dao.User, error) {
	var user dao.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return user, errors.New("无效的用户名或密码")
	}
	//验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return user, errors.New("无效的用户名或密码")
	}
	return user, nil
}

func CreatePost(post *dao.Post) error {
	if err := config.DB.Create(post).Error; err != nil {
		log.Println(err)
		return errors.New("创建文章失败")
	}
	return nil
}

func FindPost(id uint) ([]dao.Post, error) {
	var posts []dao.Post
	var result *gorm.DB
	if id == 0 {
		result = config.DB.Find(&posts)
	} else {
		result = config.DB.Where("id = ?", id).Find(&posts)
	}
	if err := result.Error; err != nil {
		log.Println(err)
		return nil, errors.New("创建文章失败")
	}
	return posts, nil
}

func UpdatePost(post dao.Post) error {
	result := config.DB.Model(&post).Select("Title", "Content", "UserId").Updates(post)
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}
	return nil
}

func DeletePost(id uint) error {
	//级联删除
	var posts []dao.Post
	if err := config.DB.Preload("Comments").Find(&posts, id).Error; err != nil {
		log.Println(err)
	}
	result := config.DB.Select("Comments").Delete(&posts, id)
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}
	return nil
}

func CreateComment(comment *dao.Comment) error {
	var count int64
	config.DB.Model(&dao.Post{}).Where("id = ?", comment.PostID).Count(&count)
	if count == 0 {
		return errors.New("文章不存在")
	}
	if err := config.DB.Create(comment).Error; err != nil {
		log.Println(err)
		return errors.New("创建评论失败")
	}
	return nil
}

func FindComment(postId uint) ([]dao.Comment, error) {
	var comments []dao.Comment
	result := config.DB.Where("post_id = ?", postId).Find(&comments).Order("created_at desc")
	if err := result.Error; err != nil {
		log.Println(err)
		return nil, result.Error
	}
	return comments, nil
}
