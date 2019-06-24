#!/usr/bin/env bash
if [ -f ~/.bash_profile ]; then
. ~/.bash_profile
fi

# beego安装
go get -u github.com/beego/bee
go get -u github.com/astaxie/beego

# MySQL
go get -u github.com/go-sql-driver/mysql
go get -u github.com/jmoiron/sqlx

# Redis
go get -u github.com/yangfei4913438/redis-full

# websocket
go get -u github.com/gorilla/websocket

# 验证码(可能需要 git 翻墙，请自行百度 git 翻墙的方法)
go get -u github.com/mojocn/base64Captcha

# 验证码不能翻墙的境内处理办法
# 方法一：
# mkdir -p $GOPATH/src/golang.org/x
# cd $GOPATH/src/golang.org/x
# git clone https://github.com/golang/image.git
#
# 方法二：
# go version > 1.11
# set env GOPROXY=https://goproxy.io