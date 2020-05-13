package main

import (
	"fmt"
	"net/http"
	"strconv"

	"helloren.cn/snippetbox/pkg/forms"
	"helloren.cn/snippetbox/pkg/models"
)

// 方法		匹配路劲			处理函数			说明
// GET 		/ 					home			显示主页
// GET 		/snippet?id=1		showSnippet 	显示片段
// GET	 	/snippet/create 	createSnippet	新的代码段表单
// POST 	/snippet/create 	createSnippet 	创建片段
// GET 		/user/signup		signupUserForm  用户注册表单
// POST 	/user/signup		signupUser		用户注册
// GET 		/user/login 		loginUserForm 	用户登录表单
// POST 	/user/login 		loginUser	 	用户登录
// POST 	/user/logout 		logoutUser	 	用户注销
// GET 		/static/ 			http.FileServer 服务特定的静态文件

//显示主页
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	//处理没有的页面 ，添加第三方路由后，这里不需要单独处理了
	// if r.URL.Path != "/" {
	// 	app.notFound(w)
	// 	return
	// }

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{Snippets: s})
}

//显示片段
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {

	//处理没有id参数
	id, err := strconv.Atoi(r.URL.Query().Get(":id")) //go标准库路由id  第三方路由:id
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

//新的代码段表单
func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

//创建片段
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	//处理不是POST请求 添加第三方路由后，这里不需要单独处理了
	// if r.Method != http.MethodPost {
	// 	w.Header().Set("Allow", "POST")
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	return
	// }

	//限制请求数据大小为4096字节
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	//填充数据，检测数据是否有错误
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//创建表单
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")
	//检测表单数据是否有错误
	if !form.Valid() {
		//有错误，重新渲染页面
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	//添加会话
	app.session.Put(r, "flash", "Snippet successfully created!")

	//添加成功重定向到显示页面
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

//用户注册表单
func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

//用户注册
func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 6)
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}

	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err == models.ErrDuplicateEmail {
		form.Errors.Add("email", "Address is already in use")
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

//用户登录表单
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

//用户登录
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Email or Password is incorrect")
		app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	// Add the ID of the current user to the session, so that they are now 'logg
	// in'.
	app.session.Put(r, "userID", id)
	// Redirect the user to the create snippet page.
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

//用户注销
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "userID")
	// Add a flash message to the session to confirm to the user that they've be
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
