package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/justinas/nosurf"
	"helloren.cn/snippetbox/pkg/models"
)

//异常恢复中间件
//注意：如果您有一个处理程序可以启动另一个goroutine（例如进行一些后台处理），然后在
//第二个goroutine将无法恢复-不能通过recoverPanic中间件来恢复，它们将导致您的应用程序退出并关闭服务器。
//需要在go func(){  if err := recover(); err != nil xxxx }() 自行处理
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

//记录日志
//运行流程:logRequest ↔ secureHeaders ↔ servemux ↔ application handler
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s %s %s\r\n", r.RemoteAddr, r.Proto, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

//安全头，防止XSS和Clickjacking攻击
//中间件位置: secureHeaders → servemux → application handler
//运行流程:secureHeaders → servemux → application handler → servemux → secureHeaders
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		//在这行之前可以处理授权认证之类的
		next.ServeHTTP(w, r)
		//在这之后的处理是无效的
		//w.WriteHeader(http.StatusForbidden)
	})
}

//请求需要认证的页面中间件
func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.authenticatedUser(r) == nil {
			http.Redirect(w, r, "/user/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//检查会话中是否存在userID值。如果不是存在，然后照常调用链中的下一个处理程序
		exists := app.session.Exists(r, "userID")
		if !exists {
			next.ServeHTTP(w, r)
			return
		}

		//从数据库中获取当前用户的详细信息。如果找不到匹配的记录，从中删除（无效的）用户ID
		//并照常调用链中的下一个处理程序。
		user, err := app.users.Get(app.session.GetInt(r, "userID"))
		if err == models.ErrNoRecord {
			app.session.Remove(r, "userID")
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}

		//请求来自有效的，通过身份验证（登录）的用户。我们创建了一个新的副本请求
		//并将用户信息添加到请求上下文中，以及调用链中的下一个处理程序使用新的请求
		ctx := context.WithValue(r.Context(), contextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

//防止跨站使用会话cookie
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}
