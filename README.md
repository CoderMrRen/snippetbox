# snippetbox
一个基于go语言开发的Snippetbox的Web应用程序 <br>
在线体验：https://www.helloren.cn/  <br>
主要功能点：包含路由，中间件，异常恢复，静态模板渲染，自定义模板函数，https，session，Token，上下文访问，防XSS，防跨站cookie会话，防彩虹表攻击等实现，麻雀虽小，五脏俱全，期待一起技术交流


## 路由
| 方法 | 匹配路径 | 处理函数 | 说明|
| :- | :------ | :---- | :---- |
|GET 	|	/ 	| home	| 显示主页 |
|GET 	| /snippet?id=1	  | showSnippet 	  | 显示片段
|GET	| /snippet/create | createSnippetForm | 片段表单
|POST 	| /snippet/create | createSnippet 	  | 创建片段
|GET 	| /user/signup	  |	signupUserForm    | 用户注册表单
|POST 	| /user/signup	  |	signupUser		  | 用户注册
|GET 	| /user/login 	  |	loginUserForm 	  | 用户登录表单
|POST 	| /user/login 	  | loginUser	 	  | 用户登录
|POST 	| /user/logout 	  | logoutUser	 	  | 用户注销
|GET 	| /static/ 		  | http.FileServer   | 静态文件

## 第三方库
| 功能点 | 地址 |
| :--- | :-- |
| 配置文件  | github.com/xiaoqidun/goini	    |
| 路由管理  | github.com/bmizerany/pat          |
| 中间件管理  | github.com/justinas/alice       |
| sessions | github.com/golangcollege/sessions |
| 防跨站点攻击 | github.com/justinas/nosurf     |
| mysql | github.com/go-sql-driver/mysql     |


## 实现参考
https://www.alexedwards.net/blog/

## 技术交流学习
QQ微信同号：136384658
