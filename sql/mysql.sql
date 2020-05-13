
--mysql数据库脚本

--创建用访问数据库的户名和密码，
--用户名：snippetbox
--密码：A4D0989DDEC964A260AAEF3FE08288E6
--支持外部网络访问：%
CREATE USER  IF NOT EXISTS 'snippetbox'@'%' IDENTIFIED BY 'A4D0989DDEC964A260AAEF3FE08288E6';

--设置snippetbox数据库权限：SELECT,INSERT
GRANT SELECT,INSERT on snippetbox.* TO 'snippetbox'@'%';

--立即生效
FLUSH PRIVILEGES;

--创建snippetbox数据库 只要数据库默认指定编码后，后续表会默认使用指定的编码创建,utf8mb4在mysql中是真正的utf-8
CREATE DATABASE IF NOT EXISTS snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

--切换数据库
use snippetbox;

--创建片段表 
--id 标题 内容 创建时间 到期时间 将列created创建索引(索引名idx_snippets_created)
CREATE TABLE IF NOT EXISTS snippets (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL,
    INDEX idx_snippets_created(created)
);

--创建用户表
CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL,
    CONSTRAINT users_uc_email UNIQUE (email) 
);