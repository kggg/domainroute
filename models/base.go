package models

import (
	"fmt"
	"time"
)

const (
	basedir        = "/store/database/dropbox/domainroute" // 程序运行的根目录
	timeLayout     = "2006-01-02 15:04:05"
	routeTablePath = "/etc/iproute2/rt_tables"
)

var (
	//basedir, _ = os.Getwd()
	iplistpath = basedir + "/iplist"
	domainpath = basedir + "/route.ini"
)

func timeConversion(t string) (int64, error) {
	times, err := time.Parse(timeLayout, t)
	if err != nil {
		return 0, fmt.Errorf("Convert tiem error:%w", err)
	}
	timeUnix := times.Unix()
	return timeUnix, nil
}
