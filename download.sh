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
