package models

import (
	"errors"
	"time"
)

var (
	//ErrNoRecord 没有片段记录
	ErrNoRecord = errors.New("models: no matching record found")
	//ErrInvalidCredentials 用户验证失败
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	//ErrDuplicateEmail 电子邮件已经存在
	ErrDuplicateEmail = errors.New("models: duplicate email")
)

//Snippet 片段对象
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

//User 用户
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}
