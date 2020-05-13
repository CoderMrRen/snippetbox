package main

import (
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	"helloren.cn/snippetbox/pkg/forms"
	"helloren.cn/snippetbox/pkg/models"
)

//模板数据对象
type templateData struct {
	AuthenticatedUser *models.User //会话认证用户
	CSRFToken         string       //防止跨站点使用cookie会话
	CurrentYear       int
	Flash             string      //会话
	Form              *forms.Form //处理表单错误
	Snippet           *models.Snippet
	Snippets          []*models.Snippet
}

//自定义模板函数
func humanDate(t time.Time) string {
	return t.Add(8 * time.Hour).Format("2006-01-02 15:04:05")
}

//自定义模板函数和函数本身对应
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	//获取dir路径下面所有扩展名.page.tmpl 页面模板
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil || pages == nil {
		return nil, fmt.Errorf("newTemplateCache filepath.Glob error:%v,%v", pages, err)
	}

	//
	for _, page := range pages {
		//获取文件名 *.page.tmpl
		name := filepath.Base(page)

		//创建一个空模板集，使用Funcs（）方法注册 自定义模板函数
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		//将页面模板文件解析为模板集
		// ts, err = template.ParseFiles(page)
		// if err != nil {
		// 	return nil, err
		// }

		//将布局模板添加到模板集
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		//将部分模板添加到模板集
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
