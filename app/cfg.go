package app

import (
	"database/sql"
	"domainroute/utils"
	"fmt"
	"os"
	"path"
	"strings"

	"gopkg.in/ini.v1"
)

type Config struct {
	Basedir         string
	Mode            string
	RouteLifeTime   int64
	RouteTablesPath string
	RouteFilePath   string
	IPListDir       string
	DBInfo
}

type DomainInfo struct {
	Domainname string
	Gateway    string
}

type DBInfo struct {
	Host   string
	User   string
	Pass   string
	Port   int
	Dname  string
	Dbconn *sql.DB
}

func NewConfig() (Config, error) {
	workdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cfg, err := ini.Load(path.Join(workdir, "./conf/app.ini"))
	if err != nil {
		panic(err)
	}
	var c Config
	c.Basedir = cfg.Section("base").Key("basedir").String()
	c.Mode = cfg.Section("base").Key("mode").String()
	c.RouteLifeTime, err = cfg.Section("base").Key("routeLifetime").Int64()
	if err != nil {
		return Config{}, fmt.Errorf("%w", err)
	}
	c.RouteTablesPath = cfg.Section("base").Key("routeTablesPath").String()
	c.RouteFilePath = cfg.Section("base").Key("routeFilePath").String()
	if strings.HasPrefix(c.RouteFilePath, "./") {
		c.RouteFilePath = path.Join(c.Basedir, c.RouteFilePath)
	}

	if c.Mode == "file" {
		c.IPListDir = cfg.Section("file").Key("iplistpath").String()
		if strings.HasPrefix(c.IPListDir, "./") {
			c.IPListDir = path.Join(c.Basedir, c.IPListDir)
		}
	}
	if c.Mode == "mysql" {
		c.DBInfo.Host = cfg.Section("mysql").Key("host").String()
		c.DBInfo.User = cfg.Section("mysql").Key("user").String()
		c.DBInfo.Pass = cfg.Section("mysql").Key("pass").String()
		c.DBInfo.Dname = cfg.Section("mysql").Key("dbname").String()

		c.DBInfo.Port, err = cfg.Section("mysql").Key("port").Int()
		if err != nil {
			fmt.Printf("Get port from config file error: %v, reset the mysql server port to default port 3306\n", err)
			c.DBInfo.Port = 3306
		}
	}
	return c, nil

}

//GetDomain 获取route.ini文件中的内容
func (c *Config) GetDomain() ([]DomainInfo, error) {
	//fmt.Println(c.RouteFilePath)
	content, err := utils.ReadFromFile(c.RouteFilePath)
	if err != nil {
		return nil, err
	}
	errA := utils.FormatParser(content)
	if errA != nil {
		return nil, errA
	}
	var domains []DomainInfo
	for _, line := range content {
		if len(line) == 0 {
			continue
		}
		domaininfo := strings.Fields(line)
		domains = append(domains, DomainInfo{domaininfo[0], domaininfo[2]})
	}
	return domains, nil
}
