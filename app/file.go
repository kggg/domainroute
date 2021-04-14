package app

import (
	"domainroute/utils"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// SaveToFile 保存域名dname解析到的IP列表iplist到文件中
func (c *Config) SaveToFile(dname string, newiplist []string) error {
	var storeFilePath string
	if strings.HasPrefix(c.IPListDir, "./") {
		storeFilePath = filepath.Join(c.Basedir, c.IPListDir, dname)
	} else {
		storeFilePath = c.IPListDir + "/" + dname
	}

	ok := utils.CheckFileExists(storeFilePath)
	if !ok {
		f, err := os.Create(storeFilePath)
		if err != nil {
			return fmt.Errorf("Create file %s error: %w", dname, err)
		}
		defer f.Close()

		var ipstring string
		now := time.Now().Format(utils.TimeLayout)

		for _, v := range newiplist {
			ipstring += v + " " + now + "\n"
		}
		f.WriteString(ipstring)

	} else {
		//从文件中读取已经保存的地址信息
		oldiplist, err := utils.ReadFromFile(storeFilePath)
		if err != nil {
			return fmt.Errorf("Read oldiplsit error: %w", err)
		}
		// 对比新解析到的地址和旧文件中保存的地址
		iplist, err := c.Compare(newiplist, oldiplist)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		var alliplist string
		for _, v := range iplist {
			alliplist += v + "\n"
		}
		//fmt.Println(filenamepath)
		err = os.WriteFile(storeFilePath, []byte(alliplist), 644)
		if err != nil {
			return fmt.Errorf("write to file %s error: %w", storeFilePath, err)
		}
	}
	return nil
}

// Compare 对比新解析得到的IP列表与已经存在文件的IP列表， 如果IP已存在，则更新此IP的存储时间， 如果没有则追加到文件的末尾
func (c *Config) Compare(newiplist, oldiplist []string) ([]string, error) {
	ipmap := make(map[string]string)
	now := time.Now().Format(utils.TimeLayout)

	for _, v1 := range newiplist {
		ipmap[v1] = now
	}

	for _, v2 := range oldiplist {

		if len(v2) == 0 {
			continue
		}

		//如果之前的保存时间超过半年没有更新， 则去掉这个IP地址
		ptime := strings.SplitN(v2, " ", 2)[1]
		ptime = strings.TrimSuffix(ptime, "\n")
		pretime, err := utils.TimeConversion(ptime)
		if err != nil {
			continue
		}
		if time.Now().Unix()-pretime >= c.RouteLifeTime {
			//同时要删除路由表中相关该地址的路由
			err := c.delroute(strings.Split(v2, " ")[0])
			if err != nil {
				return nil, err
			}
			continue
		}
		ipmap[strings.SplitN(v2, " ", 2)[0]] = ptime
	}

	var iplist []string
	for k, v := range ipmap {
		iplist = append(iplist, k+" "+v)
	}
	return iplist, nil
}

// ReadIPFormFile 从存储文件中读取IP列表
func (c *Config) ReadIPFormFile(dname string) ([]string, error) {
	filenamepath := path.Join(c.IPListDir, dname)
	content, err := utils.ReadFromFile(filenamepath)
	if err != nil {
		return nil, err
	}
	for i, v := range content {
		content[i] = strings.Split(v, " ")[0]
	}
	return content, nil
}
