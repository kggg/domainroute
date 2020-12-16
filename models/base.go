package models

import (
	"os"
)

var (
	basedir, _ = os.Getwd() // 程序运行的根目录
	iplistpath = basedir + "/iplist"
	domainpath = basedir + "/domain.txt"
)

const timeLayout = "2006-01-02 15:04:05"
