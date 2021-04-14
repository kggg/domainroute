package app

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DRouteinfo struct {
	Id         int
	Domainname string
	Ipaddr     string
	Created    time.Time
	Updated    time.Time
}

//ReadIPFromDB 从数据库中读取域名domain包含的IP地址
func (c *Config) ReadIPFromDB(domain string) ([]string, error) {
	query, err := c.findByDomainName(domain)
	if err != nil {
		return nil, err
	}
	var iplist []string
	now := time.Now().Unix()
	for _, v := range query {
		if now-v.Updated.Unix() >= c.RouteLifeTime {
			continue
		}
		iplist = append(iplist, v.Ipaddr)
	}
	return iplist, nil

}

//SaveToDB 保存解析得到的IP列表到数据库中
func (c *Config) SaveToDB(domain string, addr []string) error {
	for _, ip := range addr {
		if c.checkIPExist(domain, ip) {
			//如果此域名domain对应的ip已经有记录， 则更新更改时间
			_, err := c.DBUpdate(domain, ip)
			if err != nil {
				return err
			}
		} else {
			//插入新记录
			_, err := c.DBInsert(domain, ip)
			if err != nil {
				return err
			}
		}
	}
	return nil

}

func (c *Config) DBInsert(domain, ip string) (int64, error) {
	sqlstr := "insert into domainroute (domainname, ipaddr) values(?, ?)"
	ret, err := c.Dbconn.Exec(sqlstr, domain, ip)
	if err != nil {
		return 0, fmt.Errorf("insert db record error: %w", err)
	}
	newId, err := ret.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get insert id error: %w", err)
	}
	return newId, nil
}

func (c *Config) DBUpdate(domain, ip string) (int64, error) {
	now := time.Now()
	sqlstr := "update domainroute set updated = ? where domainname = ? and ipaddr = ?"
	res, err := c.Dbconn.Exec(sqlstr, now, domain, ip)
	if err != nil {
		return 0, fmt.Errorf("Update dbinfo error: %w", err)
	}
	n, err := res.RowsAffected() //操作影响的行数
	if err != nil {
		return 0, fmt.Errorf("get RowsAffected failed, err:%v\n", err)
	}
	return n, nil

}

func (c *Config) checkIPExist(domain, ip string) bool {
	var exist bool
	sqlstr := "select EXISTS(select 1 from domainroute where domainname = ? and ipaddr = ?)"
	res := c.Dbconn.QueryRow(sqlstr, domain, ip)
	if err := res.Scan(&exist); err != nil {
		return false
	}
	return exist
}

func (c *Config) findByDomainName(domain string) ([]DRouteinfo, error) {
	sqlstr := "select * from domainroute where domainname = ?"
	rows, err := c.Dbconn.Query(sqlstr, domain)
	if err != nil {
		return nil, fmt.Errorf("Query db error: %w", err)
	}
	var drouteinfos []DRouteinfo
	var drouteinfo DRouteinfo
	for rows.Next() {
		rows.Scan(&drouteinfo.Id, &drouteinfo.Domainname, &drouteinfo.Ipaddr, &drouteinfo.Created)
		drouteinfos = append(drouteinfos, drouteinfo)
	}
	return drouteinfos, nil
}

func (c *Config) ConnectionDB() error {
	db, err := initdb(c.User, c.Pass, c.Host, c.Dname, c.Port)
	if err != nil {
		return err
	}
	c.DBInfo.Dbconn = db
	return nil
}

func initdb(user, pass, host, dbname string, port int) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("connection db error: %w", err)
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(100 * time.Second)

	if err = db.Ping(); nil != err {
		panic("数据库链接失败: " + err.Error())
	}
	return db, nil
}
