package main

import (
	"log"
	"sync"

	"github.com/xiaoqidun/goini"
)

var once sync.Once

//Config Config
type Config struct {
	//app配置
	App struct {
		RunMode    string `goini:"run_mode"`    //程序运行模式
		AppName    string `goini:"app_name"`    //程序名称
		AppVersion string `goini:"app_version"` //程序版本
		LogPath    string `goini:"log_path"`    //日志路径
	} `goini:"app"`

	//服务器配置
	Server struct {
		EnableHTTPS bool   `goini:"enable_https"` //使用HTTPS为true，否则使用HTTP为false
		IPAddress   string `goini:"ipaddress"`    //启动服务器地址
		Port        string `goini:"port"`         //启动服务器端口
		StaticDir   string `goini:"static_dir"`   //静态文件路径
		UIDir       string `goini:"ui_dir"`
		CertFile    string `goini:"cert_file"` //HTTPS服务证书
		KeyFile     string `goini:"key_file"`  //HTTPS服务证书key
	} `goini:"server"`

	//数据库配置
	DataBase struct {
		DriverName string `goini:"driver_name"` //数据库驱动名称
		UserName   string `goini:"user_name"`   //数据库用户名
		Password   string `goini:"password"`    //数据库密码
		Network    string `goini:"network"`     //数据库连接类型
		IPAddress  string `goini:"ipaddress"`   //数据库ip地址
		Port       string `goini:"port"`        //数据库端口
		DBName     string `goini:"db_name"`     //数据库名称
		Extension  string `goini:"extension"`   //数据库名称
	} `goini:"database"`
}

//NewConfig 单例模式创建Config对象
func newConfig(path string) *Config {
	var config *Config
	once.Do(func() {
		//加载配置文件
		ini := goini.NewGoINI()
		if err := ini.LoadFile(path); err != nil {
			log.Fatalf("%v\r\n", err)
			return
		}
		config = &Config{}
		if err := ini.MapToStruct(config); err != nil {
			log.Fatalf("%v\r\n", err)
			return
		}
	})
	return config
}
