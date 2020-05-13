package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/justinas/nosurf"
	"helloren.cn/snippetbox/pkg/models"
)

//服务器内部错误，记录错误消息和堆栈跟踪,向客户端发送500内部错误状态码
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s\r\n", err.Error(), debug.Stack())
	log.Printf(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

//客户端错误，不记录日志，向客户端发送指定状态码
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

//访问不存在页面，不记录日志
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

//默认值 页脚展示年份
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	//token
	td.CSRFToken = nosurf.Token(r)
	//会话用户
	td.AuthenticatedUser = app.authenticatedUser(r)
	//当前年份
	td.CurrentYear = time.Now().Year()
	//session 弹出一次性即时消息
	td.Flash = app.session.PopString(r, "flash")
	return td
}

//渲染模板
func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {

	ts, ok := app.templataCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
	}

	//执行渲染模板
	buf := new(bytes.Buffer)
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
	}

	buf.WriteTo(w)
}

//验证用户会话，如果请求来自未经身份验证的用户,返回nil,否则返回当前用户信息
func (app *application) authenticatedUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(contextKeyUser).(*models.User)
	if !ok {
		return nil
	}
	return user
}
