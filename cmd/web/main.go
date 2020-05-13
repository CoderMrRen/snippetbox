package main

//
//go官方开发者英文站：go.dev
//go官方开发者中文站 golangclub.com

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql" // New import
)

func init() {
	debugppfof()
}

func main() {

	//创建app
	app := newApplicaton(newConfig("../../config.ini"))

	//初始化服务器，
	srv := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", app.config.Server.IPAddress, app.config.Server.Port),
		ErrorLog:       log.New(log.Writer(), log.Prefix(), log.Flags()), //记录错误日志
		Handler:        app.routes(),
		TLSConfig:      newTLSConfig(),
		IdleTimeout:    time.Minute,     //客户端空闲时间1分钟后，关闭与客户端连接
		ReadTimeout:    5 * time.Second, //值越小可以防止洪水攻击
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: (1024*1024)/2 - 4096, //客户端最大请求头，默认1M,这里设置为0.5M,go会默认添加4096字节，所有再减去4096，达到精确控制
	}

	//main函数退出 关闭数据库连接
	defer app.snippets.DB.Close()

	log.Printf("Starting server on %s", srv.Addr)

	//http重定向到https
	go http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("https://%s%s", r.Host, r.URL.String())
		log.Printf(url)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}))

	//启动服务器
	var err error
	if !app.config.Server.EnableHTTPS {
		err = srv.ListenAndServe()
	} else {
		err = srv.ListenAndServeTLS(app.config.Server.CertFile, app.config.Server.KeyFile)
	}

	if err != nil {
		log.Panicf("ListenAndServe err:%v\r\n", err)
	}
}

//newTLSConfig tls配置
func newTLSConfig() *tls.Config {
	tlsConfig := &tls.Config{
		//tls 将使用服务器端CipherSuites列表中顺序的密码套件，而不是选中客户端请求过来的密码套件
		PreferServerCipherSuites: true,
		//TLS握手期间应选择握手类型，其它默认方式会占用很高的CPU,tls.X25519, tls.CurveP256是不会占高CPU,有助于提高服务器性能
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},

		//tls 使用密码套件
		CipherSuites: []uint16{
			//tls 1.2
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		MinVersion: tls.VersionTLS12, //tls 版本，最低高版本TLS 1.0
		MaxVersion: tls.VersionTLS12, //最高版本TLS 1.2
	}
	return tlsConfig
}

//ppfof http://localhost:6060/debug/pprof/
func debugppfof() {
	go func() {
		//ip := "0.0.0.0:6060"
		ip := "localhost:6060" //不会暴漏在公网上
		if err := http.ListenAndServe(ip, nil); err != nil {
			fmt.Printf("start pprof failed on %s\n", ip)
			os.Exit(1)
		}
	}()
}
