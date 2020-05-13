package main

import (
	"net/http"

	pat "github.com/bmizerany/pat"    //第三发插件 路由管理
	alice "github.com/justinas/alice" //第三发插件 中间件管理
)

func (app *application) routes() http.Handler {

	//异常恢复---->记录日志 --->处理安全头---->http handler
	//运行流程:recoverPanic ↔ logRequest ↔ secureHeaders ↔ servemux ↔ application handler
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	//标准库服务器路由
	// mux := http.NewServeMux()
	// mux.HandleFunc("/", app.home)
	// mux.HandleFunc("/snippet", app.showSnippet)
	// mux.HandleFunc("/snippet/create", app.createSnippet)
	// //处理css,图片，js资源等
	// mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir(app.Config.StaticDir))))

	//动态路由，用于session管理，静态文件路径static下的不要session
	//noSurf,防止跨域名使用cookie会话
	dynmicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)

	//第三方路由
	mux := pat.New()
	//片段路由
	mux.Get("/", dynmicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynmicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynmicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynmicMiddleware.ThenFunc(app.showSnippet))

	//用户路由
	mux.Get("/user/signup", dynmicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynmicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynmicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynmicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynmicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.logoutUser))

	mux.Get("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(app.config.Server.StaticDir))))

	return standardMiddleware.Then(mux)
}
