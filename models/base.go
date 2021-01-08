package models

const (
	basedir        = "/store/database/dropbox/domainroute" // 程序运行的根目录
	timeLayout     = "2006-01-02 15:04:05"
	routetablePath = "/etc/iproute2/rt_tables"
)

var (
	//basedir, _ = os.Getwd()
	iplistpath = basedir + "/iplist"
	domainpath = basedir + "/route.ini"
)
