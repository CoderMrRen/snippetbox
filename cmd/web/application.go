package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golangcollege/sessions"
	"helloren.cn/snippetbox/pkg/models/mysql"
)

type contextKey string

var contextKeyUser = contextKey("user")

type application struct {
	config        *Config
	session       *sessions.Session             //会话session
	snippets      *mysql.SnippetModel           //片段模型
	templataCache map[string]*template.Template //页面模板缓存
	users         *mysql.UserModel              //用户模型
}

//newApplication
func newApplicaton(config *Config) *application {
	//初始化日志 设置全局log
	if config.App.RunMode == "release" {
		//日志文件
		f, err := os.OpenFile(config.App.LogPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			log.Panicf("neApplication os.OpenFile err:%v\r\n", err)
		}
		log.SetOutput(f)
	} else {
		//标准输出
		log.SetOutput(os.Stdout)
	}
	//设置log样式
	log.SetPrefix("[INFO]\t")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//初始化模板
	templateCache, err := newTemplateCache(config.Server.UIDir)
	if err != nil {
		log.Fatalf("newTemplateCache error:%v\r\n", err)
	}

	//连接数据库
	db, err := openDB(config)
	if err != nil {
		log.Fatalf("openDB fail:%s\r\n", err)
	}

	//会话session 密钥
	secret := "6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge"
	session := sessions.New([]byte(secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true                      //设置安全cookie会话
	session.SameSite = http.SameSiteStrictMode //阻止用户浏览器中跨站点发送的会话cookie使用的情况，但是只有71%的浏览器支持。通过nosurf库来解决这个问题
	session.Persist = false                    //持续设置会话cookie是否应为持久性（即是否应在用户关闭浏览器后保留它）

	return &application{
		config:        config,
		session:       session,
		snippets:      &mysql.SnippetModel{DB: db},
		templataCache: templateCache,
		users:         &mysql.UserModel{DB: db},
	}
}

//打开数据库
func openDB(config *Config) (*sql.DB, error) {
	//name:pwd@tcp(localhost:3306)/snippetbox?parseTime=true
	conn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s%s", config.DataBase.UserName, config.DataBase.Password, config.DataBase.Network, config.DataBase.IPAddress, config.DataBase.Port, config.DataBase.DBName, config.DataBase.Extension)
	db, err := sql.Open(config.DataBase.DriverName, conn)
	if err != nil {
		return nil, err
	}

	//设置并发打开连接的最大数量
	//设置为小于或等于0表示没有最大限制
	//如果最大已达到打开的连接数，并且需要一个新的连接，Go将等待直到其中一个连接释放并变为空闲。
	//从一个用户角度，这意味着他们的HTTP请求将一直挂起，直到建立连接为止
	//db.SetMaxOpenConns(25)

	//设置连接池中最大空闲连接数。
	//小于或等于0将表示不保留任何空闲连接。
	db.SetMaxIdleConns(25)

	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
